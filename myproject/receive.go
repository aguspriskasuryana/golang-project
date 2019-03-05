package main

import (
    "encoding/json"
    "errors"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/streadway/amqp"
    "gopkg.in/olivere/elastic.v3"
    "log"
    "time"
)

const (
    indexName    = "myapplications"
    docType      = "log"
    appName      = "myApp"
    indexMapping = `{
                        "mappings" : {
                            "log" : {
                                "properties" : {
                                    "id" : { "type" : "string", "index" : "not_analyzed" },
                                    "created" : { "type" : "string" }
                                }
                            }
                        }
                    }`
)

type NewsReceive struct {
    Id      int    `form:"id" json:"id"`
    Author  string `form:"author" json:"author"`
    Body    string `form:"body" json:"body"`
    Created string `form:"created" json:"created"`
}

type LogNewsReceive struct {
    Id      string `json:"id"`
    Created string `json:"created"`
}

func failOnErrorReceive(err error, msgerr string, msgsuc string) {
    if err != nil {
        log.Fatalf("%s: %s", msgerr, err)

    } else {
        fmt.Printf("%s\n", msgsuc)
    }

}

func savedata(newsx string, client *elastic.Client) {
    //testing langsung ke db tanpa message qeue
    db := Conn()
    news2 := NewsReceive{}
    json.Unmarshal([]byte(newsx), &news2)

    id := news2.Id
    author := news2.Author
    body := news2.Body
    created := news2.Created
    insForm, err := db.Prepare("INSERT INTO News(id, author, body, created) VALUES(?,?,?,?)")
    if err != nil {
        panic(err.Error())
    }
    insForm.Exec(id, author, body, created)

    sid := fmt.Sprintf("%d", id)
    err = createIndexWithLogsIfDoesNotExist(client, sid, created)
    if err != nil {
        panic(err)
    }

    //mainelastic(news2)
    //end testing db
}

func createIndexWithLogsIfDoesNotExist(client *elastic.Client, id string, created string) error {
    exists, err := client.IndexExists(indexName).Do()
    if err != nil {
        return err
    }

    if exists {
        //return nil
        addNew(client, id, created)

    } else {
        res, err2 := client.CreateIndex(indexName).
            Body(indexMapping).
            Do()

        if err2 != nil {
            return err2
        }
        if !res.Acknowledged {
            return errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
        }

        addNew(client, id, created)
    }

    fmt.Printf("masuk")

    return nil
}

func addNew(client *elastic.Client, id string, created string) error {

    l := LogNewsReceive{
        Id:      id,
        Created: created,
    }

    _, err1 := client.Index().
        Index(indexName).
        Type(docType).
        BodyJson(l).
        Do()

    if err1 != nil {
        return err1
    }
    return nil
}

func mainreceive() {

    client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
    if err != nil {
        panic(err)
    }

    fmt.Println("Connecting to RabbitMQ")
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") //Insert the  connection string
    failOnErrorReceive(err, "Failed to connect to RabbitMQ server", "Successfully connected to RabbitMQ server")
    defer conn.Close()
    failOnErrorReceive(err, "RabbitMQ connection failure", "RabbitMQ Connection Established")

    ch, err := conn.Channel()
    failOnErrorReceive(err, "Failed to open a channel", "Opened the channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "news", //name
        //"ha.monitoring",
        true,
        false, //delete when unused
        false, //exclusive
        false, //no-wait
        nil,   //arguements
    )
    failOnErrorReceive(err, "Failed to declare the queue", "Declared the queue")

    msgs, err := ch.Consume(
        q.Name, //queue
        "",     //consumer
        true,   //auto-ack
        false,  //exclusive
        false,  //no-local
        false,  //no-wait
        nil,    //args
    )
    failOnErrorReceive(err, "Failed to register a consumer ", "Registered the consumer")

    msgCount := 0
    go func() {
        for d := range msgs {

            msgCount++

            fmt.Printf("\nMessage Count: %d, Message Body: %s\n", msgCount, d.Body)
            //fmt.Printf(string(d.Body))
            savedata(string(d.Body), client)
        }
    }()

    select {
    case <-time.After(time.Second * 10):
        fmt.Printf("Total Messages Fetched: %d\n", msgCount)
        fmt.Println("No more messages in queue. Timing out...")

    }

}
