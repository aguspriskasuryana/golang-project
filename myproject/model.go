package main

import (
	//"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"time"
)

//News struct
type News struct {
	Id      int    `form:"id" json:"id"`
	Author  string `form:"author" json:"author"`
	Body    string `form:"body" json:"body"`
	Created string `form:"created" json:"created"`
}

var newone []News

//Response struct
type ResponseNew struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []News
}

// funtion untuk memparsing data MySQL ke JSON
func returnAllNews(w http.ResponseWriter, r *http.Request) {
	var news News               // variable untuk memetakan data yang terbagi menjadi 3 field
	var arr_news []News         // menampung variable ke dalam bentuk slice
	var responseNew ResponseNew //variable untuk menampung data arr yang nantinya akan diubah menjadi bentuk JSON
	db := Conn()

	rows, err := db.Query("Select id,author,body,created from news ORDER BY created DESC")
	if err != nil {
		log.Print(err)
	}
	responseNew.Status = 1 //mengisi valus status = 1
	responseNew.Message = "Success"

	//bentuk perulangan untuk me render data dari mySQL ke struct dan slice data
	for rows.Next() {
		if err := rows.Scan(&news.Id, &news.Author, &news.Body, &news.Created); err != nil {
			log.Fatal(err.Error())
			responseNew.Status = 0
			responseNew.Message = "Failed"
		} else {
			arr_news = append(arr_news, news)
		}
	}
	responseNew.Data = arr_news            // mengisi komponen Data dengan data slice arr
	json.NewEncoder(w).Encode(responseNew) //mengubah data struct menjadi JSON

}

func CreateNewEndpoint(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(req)
	r.ParseForm()
	if r.Method == "POST" {
		var news News
		t := time.Now()
		_ = json.NewDecoder(r.Body).Decode(&news)
		//ids, _ := t.Format("150405")

		ids, _ := strconv.Atoi(t.Format("150405"))
		news.Id = ids //atau dibuatkan generate unix id
		news.Author = r.FormValue("author")
		news.Body = r.FormValue("body")
		news.Created = t.Format("2006-01-02 15:04:05")

		var newonex []News
		newonex = append(newonex, news)
		n, _ := json.Marshal(news)
		// Convert bytes to string.
		newstring := string(n)
		//json.NewEncoder(w).Encode(newone)

		Send(newstring)
		New(w, r)

	}
}

func GetNewEndpoint(w http.ResponseWriter, r *http.Request) {

	mainreceive()

}
func FetchNews(w http.ResponseWriter, r *http.Request) {

	mainlist(w, r)

}
