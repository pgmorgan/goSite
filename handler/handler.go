package handler

import (
	"net/http"

	"github.com/pgmorgan/goSite/db"
	"github.com/pgmorgan/goSite/tpl"
)

func Index(w http.ResponseWriter, req *http.Request) {
	data := struct {
		books []db.Books
		//USER INFO
	}{
		books: list,
	}

	list, err := db.DBlist()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
		return
	}

	tpl.TPL.ExecuteTemplate(w, "books.gohtml", data)
}
