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
	list, err := db.DBlist()
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
		users.AlreadyLoggedIn(req),
	}

	tpl.TPL.ExecuteTemplate(w, "index.gohtml", data)
}

func Insert(w http.ResponseWriter, req *http.Request) {
	book := db.Book{
		Title:  req.FormValue("title"),
		Author: req.FormValue("author"),
		Price:  req.FormValue("price"),
		// Currency: req.FormValue("currency"),
		BuyLink: req.FormValue("buylink"),
		ID:      req.FormValue("id"),
	}
	err := db.DBinsertOne(book)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, req *http.Request) {
	title, err := url.QueryUnescape(req.FormValue("urltitle"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+
			err.Error(), http.StatusInternalServerError)
	}
	err = db.DBdeleteOne(title)

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
	id := req.FormValue("id")
	alreadyListed, err := db.DBidAlreadyListed(id)
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
