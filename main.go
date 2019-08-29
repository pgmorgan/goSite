package main

import (
	"net/http"

	"github.com/pgmorgan/goSite/handler"
)

func main() {
	http.HandleFunc("/add", handler.Add)
	http.HandleFunc("/search", handler.Search)
	http.HandleFunc("/insert", handler.Insert)
	http.HandleFunc("/delete", handler.Delete)
	http.HandleFunc("/", handler.Index)
	http.ListenAndServe(":80", nil)
	// handler.Index()
}
