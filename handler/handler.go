package handler

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pgmorgan/goSite/bookapi"
	"github.com/pgmorgan/goSite/db"
	"github.com/pgmorgan/goSite/tpl"
	"github.com/pgmorgan/goSite/users"
)

/*	INDEX HANDLER - SERVES INDEX PAGE	*/
func Index(w http.ResponseWriter, req *http.Request) {
	/*	AlreadyLoggedIn() method checks the browser cookie which
	**	is linked to a user account in a map[string]string{}.
	**	If the cookie and user are stored in the map it returns
	**	true, email
	 */
	userEmail, loggedIn := users.AlreadyLoggedIn(req)
	var err error
	var list []db.Book
	/*	If not logged in, you can insert/remove books from a public account	*/
	if loggedIn {
		list, err = db.DBlist(userEmail)
	} else {
		list, err = db.DBlist("public")
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	data := struct {
		Dlist     []db.Book
		DloggedIn bool
	}{
		list,
		loggedIn,
	}

	tpl.TPL.ExecuteTemplate(w, "index.gohtml", data)
}

/*	INSERT BOOK HANDLER - Not exported.  Called from Add() handler below	*/
func insert(w http.ResponseWriter, req *http.Request, book db.Book) {
	var err error

	userEmail, loggedIn := users.AlreadyLoggedIn(req)
	if loggedIn {
		err = db.DBinsertOne(book, userEmail)
	} else {
		err = db.DBinsertOne(book, "public")
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

/*	DELETE BOOK HANDLER	*/
func Delete(w http.ResponseWriter, req *http.Request) {
	var err error

	title, err := url.QueryUnescape(req.FormValue("urltitle"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
	}
	userEmail, loggedIn := users.AlreadyLoggedIn(req)
	if loggedIn {
		err = db.DBdeleteOne(title, userEmail)
	} else {
		err = db.DBdeleteOne(title, "public")
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

/*	SEARCH FOR TOP TEN BOOK MATCHES HANDLER	*/
func Search(w http.ResponseWriter, req *http.Request) {
	_, loggedIn := users.AlreadyLoggedIn(req)
	title := req.FormValue("title")
	results, err := bookapi.FindTopTen(title)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		Results   bookapi.JsonObjects
		DloggedIn bool
	}{
		results,
		loggedIn,
	}

	tpl.TPL.ExecuteTemplate(w, "searchResults.gohtml", data)
}

/*	ADD BOOK HANDLER - Checks if book is already in the collection and calls insert() method above */
func Add(w http.ResponseWriter, req *http.Request) {
	var err error
	var alreadyListed bool

	id := req.FormValue("id")
	userEmail, loggedIn := users.AlreadyLoggedIn(req)
	if loggedIn {
		alreadyListed, err = db.DBidAlreadyListed(id, userEmail)
	} else {
		alreadyListed, err = db.DBidAlreadyListed(id, "public")
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
	}
	if alreadyListed {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	result, err := bookapi.FindOne(id)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
	}
	book := db.Book{
		Title:    result.VolumeInfo.Title,
		Author:   result.VolumeInfo.Author[0],
		ThumbURL: result.VolumeInfo.ImgLink.Thumb,
		ID:       result.ID,
	}
	if result.SaleInfo.RetailPrice.Amount != 0 {
		book.Price = strconv.FormatFloat(result.SaleInfo.RetailPrice.Amount, 'f', 2, 64)
		book.BuyLink = result.SaleInfo.BuyLink
	}
	insert(w, req, book)
}
