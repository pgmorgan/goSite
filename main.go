package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pgmorgan/goSite/handler"
	"github.com/pgmorgan/goSite/users"
)

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatalln("Missing PORT environment variable in .env file at root of repository")
	}
	port = ":" + port

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/add", handler.Add)
	http.HandleFunc("/search", handler.Search)
	http.HandleFunc("/delete", handler.Delete)
	http.HandleFunc("/signup", users.SignUp)
	http.HandleFunc("/logout", users.Logout)
	http.HandleFunc("/login", users.Login)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", handler.Index)
	http.ListenAndServe(port, nil)
}
