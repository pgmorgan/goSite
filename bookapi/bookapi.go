package bookapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pgmorgan/goSite/db"
)

func FindTopTen(title string) ([]db.Book, error) {
	// url := "https://www.googleapis.com/books/v1/volumes?q=harry+potter"
	url := "http://api.open-notify.org/astros.json"

	type rawJSON struct {
		Number int `json:"number"`
	}

	clientTimeout := http.Client{
		Timeout: time.Second * 4,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "book-request-api-call")

	res, err := clientTimeout.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bookInfo := rawJSON{}
	err = json.Unmarshal(body, &bookInfo)
	if err != nil {
		return nil, err
	}

	fmt.Println(bookInfo.Number)
	return nil, nil
}
