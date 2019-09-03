package handler

import (
	"log"
	"net/http"
	"net/url"

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
	var uri string

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

	// Currency: req.FormValue("currency"),

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
		book = db.Book{
			Price:   result.SaleInfo.RetailPrice.Amount,
			BuyLink: result.SaleInfo.BuyLink,
		}
	}
	insert(w, req, book)
}

// fmt.Println("Reached here!")
// http.Redirect(w, req, "/insert"+
// 	"?title="+result.VolumeInfo.Title+
// 	"&author="+result.VolumeInfo.Author[0]+
// 	"&thumb="+url.PathEscape(result.VolumeInfo.ImgLink.Thumb),
// 	http.StatusSeeOther)
