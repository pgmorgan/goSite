package main

import (
	"github.com/pgmorgan/goSite/handler"
)

func main() {
	// http.HandleFunc("/", handler.Index)
	// http.ListenAndServe(":8080", nil)
	handler.Index()
}
