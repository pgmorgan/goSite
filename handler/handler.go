package handler

import (
	"net/http"

	"github.com/pgmorgan/go2/046_mongodb/15_postgres/config"
	"github.com/pgmorgan/goSite/db"
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

	config.TPL.ExecuteTemplate(w, "books.gohtml", data)
}
