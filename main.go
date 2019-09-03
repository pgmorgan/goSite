package main

import (
	"net/http"

	"github.com/pgmorgan/goSite/handler"
	"github.com/pgmorgan/goSite/users"
)

func main() {
	http.HandleFunc("/add", handler.Add)
	http.HandleFunc("/search", handler.Search)
	http.HandleFunc("/insert", handler.Insert)
	http.HandleFunc("/delete", handler.Delete)
	http.HandleFunc("/signup", users.SignUp)
	http.HandleFunc("/logout", users.Logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", handler.Index)
	http.ListenAndServe(":8080", nil)
}
