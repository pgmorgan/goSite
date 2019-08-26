package main

import (
	"net/http"

	"github.com/pgmorgan/goSite/handler"
)

func main() {
	http.HandleFunc("/insert", handler.Insert)
	http.HandleFunc("/", handler.Index)
	http.ListenAndServe(":8080", nil)
	// handler.Index()
}
