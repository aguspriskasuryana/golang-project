package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

var tmpl = template.Must(template.ParseGlob("form/*"))

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

//fungsi main
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/new", New) // untuk testing sementara
	router.HandleFunc("/producer", CreateNewEndpoint)
	router.HandleFunc("/consumer", GetNewEndpoint)
	router.HandleFunc("/getnewspagination", FetchNews).Methods("GET")
	router.HandleFunc("/getnews", returnAllNews).Methods("GET") // menjalurkan URL untuk dapat mengkases data JSON API new
	//router.HandleFunc("/insert2/{id}", CreatePersonEndpoint).Methods("POST")
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", router))

}
