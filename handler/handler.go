package handler

import (
	"log"
	"net/http"

	"github.com/pgmorgan/goSite/db"
	"github.com/pgmorgan/goSite/tpl"
)

func Index(w http.ResponseWriter, req *http.Request) {
	// data := struct {
	// 	books []db.Book
	// }{
	// 	books: list,
	// }

	list, err := db.PublicList()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	tpl.TPL.ExecuteTemplate(w, "index.gohtml", list)
}

func Insert(w http.ResponseWriter, req *http.Request) {
	book := db.Book{
		Title:	"Life of Pi",
	}
	err := db.PublicInsertOne(book)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	Index(w, req)
}
