package main

import (
	"fmt"
	"gopkg.in/olivere/elastic.v3"
	"net/http"
	"reflect"
	"strconv"
)

//func mainlistx(w http.ResponseWriter, r *http.Request) {
//	tmpl.ExecuteTemplate(w, "New", nil)
//}

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
	var authorMatch = make(map[int]string)
	var bodyMatch = make(map[int]string)
	for selDB.Next() {
		var id int
		var author, body, created string
		err1 = selDB.Scan(&id, &author, &body, &created)
		if err1 != nil {
			panic(err1.Error())
		}
		authorMatch[id] = author
		bodyMatch[id] = body

	}

	termQuery := elastic.MatchAllQuery{}

	res, err := client.Search(indexName).
		Index(indexName).
		Query(termQuery).
		Sort("time", true).
		Do()

	if err != nil {
		return err
	}

	fmt.Println("Logs found:")

	nrfromelastic := NewsReceive{}
	resxfromelastic := []NewsReceive{}
	var l LogNewsReceive
	for _, item := range res.Each(reflect.TypeOf(l)) {
		l := item.(LogNewsReceive)

		idint, _ := strconv.Atoi(l.Id)

		nrfromelastic.Id = idint
		nrfromelastic.Author = authorMatch[idint]
		nrfromelastic.Body = bodyMatch[idint]
		nrfromelastic.Created = l.Created
		resxfromelastic = append(resxfromelastic, nrfromelastic)
		fmt.Printf("created: %s ID: %s\n", l.Created, l.Id)
	}

	tmpl.ExecuteTemplate(w, "Index", resxfromelastic)
	return nil
}
