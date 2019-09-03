package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pgmorgan/goSite/bookapi"
	"github.com/pgmorgan/goSite/db"
	"github.com/pgmorgan/goSite/tpl"
	"github.com/pgmorgan/goSite/users"
)

func Index(w http.ResponseWriter, req *http.Request) {
	userEmail, loggedIn := users.AlreadyLoggedIn(req)
	var err error
	var list []db.Book
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

func Insert(w http.ResponseWriter, req *http.Request) {
	var err error

	book := db.Book{
		Title:  req.FormValue("title"),
		Author: req.FormValue("author"),
		Price:  req.FormValue("price"),
		// Currency: req.FormValue("currency"),
		BuyLink: req.FormValue("buylink"),
		ID:      req.FormValue("id"),
	}
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

func Search(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	results, err := bookapi.FindTopTen(title)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
	}
	tpl.TPL.ExecuteTemplate(w, "searchResults.gohtml", results)
}

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
	if result.SaleInfo.RetailPrice.Amount == 0 {
		fmt.Println("Reached here!")
		http.Redirect(w, req, "/insert"+
			"?title="+result.VolumeInfo.Title+
			"&author="+result.VolumeInfo.Author[0],
			http.StatusSeeOther)
	} else {
		http.Redirect(w, req, "/insert"+
			"?title="+result.VolumeInfo.Title+
			"&author="+result.VolumeInfo.Author[0]+
			"&id="+result.ID+
			"&price="+strconv.FormatFloat(result.SaleInfo.RetailPrice.Amount, 'f', 2, 64)+
			// "&currency="+result.SaleInfo.RetailPrice.Currency+
			"&buylink="+result.SaleInfo.BuyLink,
			http.StatusSeeOther)
	}
}
