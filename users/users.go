package users

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pgmorgan/goSite/tpl"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Email    string
	Password []byte
}

var dbSessions = map[string]string{} // session ID, user ID
var dbUsers = map[string]user{}      // user ID, user

func init() {
	populateDbUsers()
}

func populateDbUsers() {
	f, _ := os.Open("./user")
	defer f.Close()
	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = ';'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		dbUsers[record[0]] = user{record[0], []byte(record[1])}
	}
}

func AlreadyLoggedIn(req *http.Request) (string, bool) {
	c, err := req.Cookie("session")
	if err != nil {
		return "", false
	}
	em := dbSessions[c.Value]
	_, ok := dbUsers[em]
	return em, ok
}

func SignUp(w http.ResponseWriter, req *http.Request) {
	if _, ok := AlreadyLoggedIn(req); ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u user
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		em := req.FormValue("email")
		p := req.FormValue("password")

		// username taken?
		if _, ok := dbUsers[em]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = em
		// store user in dbUsers
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		u = user{em, bs}
		dbUsers[em] = u
		csv := append([]byte(u.Email), ';')
		csv = append(csv, u.Password...)
		csv = append(csv, '\n')

		f, err := os.OpenFile("./user", os.O_APPEND|os.O_WRONLY, 0600)
		check(err)
		defer f.Close()
		_, err = f.Write(csv)
		check(err)

		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.TPL.ExecuteTemplate(w, "signup.gohtml", nil)
}

func Logout(w http.ResponseWriter, req *http.Request) {
	if _, ok := AlreadyLoggedIn(req); !ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie("session")
	// delete the session
	delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Login(w http.ResponseWriter, req *http.Request) {
	if _, ok := AlreadyLoggedIn(req); ok {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u user
	// process form submission
	if req.Method == http.MethodPost {
		em := req.FormValue("email")
		p := req.FormValue("password")
		u, ok := dbUsers[em]
		if !ok {
			http.Error(w, "<h2>Username and/or password do not match</h2>", http.StatusForbidden)
			return
		}
		// does the entered password match the stored password?
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = em
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.TPL.ExecuteTemplate(w, "login.gohtml", u)
}
