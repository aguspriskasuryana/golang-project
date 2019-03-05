package main

import (
    //"encoding/json"
    "fmt"
    "github.com/streadway/amqp"
    "log"
)

func failOnError(err error, msgerr string, msgsuc string) {
    if err != nil {
        log.Fatalf("%s: %s", msgerr, err)

    } else {
        fmt.Printf("%s\n", msgsuc)
    }

}

func Send(news string) {
    //newone
    //fmt.Println(newone)
    // Connect to RabbitMQ server

    fmt.Println("Connecting to RabbitMQ ...")
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") //Insert the  connection string
    failOnError(err, "RabbitMQ connection failure", "RabbitMQ Connection Established")
    defer conn.Close()

    //Connect to the channel

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel", "Opened the channel")
    defer ch.Close()

    //Declare the queue where messages need to be sent. Queue will be created if not already there
    q, err := ch.QueueDeclare(
        "news", //name
        true,   //durable
        false,  //delete when unused
        false,  //exclusive
        false,  //no-wait
        nil,    //arguements
    )

    failOnError(err, "Failed to declare the queue", "Declared the queue")

    body := news

    //Publish to the queue

    err = ch.Publish(
        "",     //exchange
        q.Name, //routing key
        false,  //mandatory
        false,  //immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body),
        })

    failOnError(err, "Failed to publish a message ", "Published the message")

}
