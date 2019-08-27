package handler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/pgmorgan/goSite/bookapi"
	"github.com/pgmorgan/goSite/db"
	"github.com/pgmorgan/goSite/tpl"
)

func Index(w http.ResponseWriter, req *http.Request) {
	// data := struct {
	// 	books []db.Book
	// }{
	// 	books: list,
	// }

	list, err := db.DBlist()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	tpl.TPL.ExecuteTemplate(w, "index.gohtml", list)
}

func Insert(w http.ResponseWriter, req *http.Request) {
	book := db.Book{
		Title: req.FormValue("title"),
	}
	err := db.DBinsertOne(book)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, req *http.Request) {
	title, err := url.QueryUnescape(req.FormValue("urltitle"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
	}
	err = db.DBdeleteOne(title)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func Search(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	results, err := bookapi.FindTopTen(title)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
	}
	tpl.TPL.ExecuteTemplate(w, "index2.gohtml", results)
}
