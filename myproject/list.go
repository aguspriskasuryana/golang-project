package main

import (
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"net/http"
	"reflect"
)

func mainlist(w http.ResponseWriter, r *http.Request) {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		panic(err)
	}

	err = findAndPrintAppLogs(client, w, r)
	if err != nil {
		panic(err)
	}
}

func findAndPrintAppLogs(client *elastic.Client, w http.ResponseWriter, r *http.Request) error {
	//get all new from db
	//returnAllNews(w,r)
	db := Conn()
	selDB, err1 := db.Query("SELECT * FROM news ORDER BY created DESC")
	if err1 != nil {
		panic(err1.Error())
	}
	nr := NewsReceive{}
	resx := []NewsReceive{}
	for selDB.Next() {
		var id int
		var author, body, created string
		err1 = selDB.Scan(&id, &author, &body, &created)
		if err1 != nil {
			panic(err1.Error())
		}
		nr.Id = id
		nr.Author = author
		nr.Body = body
		nr.Created = created
		resx = append(resx, nr)
	}

	termQuery := elastic.MatchAllQuery{}

	res, err := client.Search(indexName).
		Index(indexName).
		Query(termQuery).
		Sort("created", false).
		Do()

	if err != nil {
		return err
	}

	fmt.Println("Logs found:")
	var l LogNewsReceive
	for _, item := range res.Each(reflect.TypeOf(l)) {
		l := item.(LogNewsReceive)
		fmt.Printf("created: %s ID: %s\n", l.Created, l.Id)
	}

	tmpl.ExecuteTemplate(w, "Index", resx)
	return nil
}
