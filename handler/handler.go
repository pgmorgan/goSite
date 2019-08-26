package handler

import (
	"github.com/pgmorgan/goSite/db"
)

// func Index(w http.ResponseWriter, req *http.Request) {
func Index() {
	// data := struct {
	// 	books []db.Book
	// 	//USER INFO
	// }{
	// 	books: list,
	// }

	// list, err := db.DBlist()
	db.Launch()
	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError)+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// tpl.TPL.ExecuteTemplate(w, "books.gohtml", nil)
}
