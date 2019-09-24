package bookapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

/*	STRUCTS TO READ JSON FROM GOOGLE BOOKS API	*/

type Volume struct {
	Title   string     `json:"title"`
	Author  []string   `json:"authors"`
	ImgLink imgLinkObj `json:"imageLinks"`
}

type imgLinkObj struct {
	Thumb string `json:"thumbnail"`
}

type ItemInfo struct {
	ID         string   `json:"id"`
	VolumeInfo Volume   `json:volumeInfo"`
	Author     []string `json:"authors"`
}

type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currencyCode"`
}

type Sale struct {
	RetailPrice Price  `json:"retailPrice"`
	BuyLink     string `json:"buyLink"`
}

type JsonObjects struct {
	Items []ItemInfo `json:"items"`
}

type JsonObject struct {
	VolumeInfo Volume `json:"volumeInfo"`
	SaleInfo   Sale   `json:"saleInfo"`
	ID         string `json:"id"`
}

/*	METHOD TO FIND TOP TEN MATCHES FROM GOOGLE BOOK API	*/
func FindTopTen(title string) (JsonObjects, error) {
	obj := JsonObjects{}
	apiKey, exists := os.LookupEnv("GOOGLE_DEV_API_KEY")
	if !exists {
		return obj, errors.New(`Missing GOOGLE_DEV_API_KEY environment
								 variable in .env file at root of repository.`)
	}
	uri := "https://www.googleapis.com/books/v1/volumes?q=" +
		url.QueryEscape(title) + "&key=" + apiKey

	clientTimeout := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return obj, err
	}
	req.Header.Set("User-Agent", "book-request-api-call")
	res, err := clientTimeout.Do(req)
	if err != nil {
		return obj, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

/*	METHOD TO FETCH ONE BOOK BY ID ONCE CHOSEN BY USER	*/
func FindOne(id string) (JsonObject, error) {
	obj := JsonObject{}
	apiKey, exists := os.LookupEnv("GOOGLE_DEV_API_KEY")
	if !exists {
		return obj, errors.New(`Missing GOOGLE_DEV_API_KEY environment
								 variable in .env file at root of repository.`)
	}
	uri := "https://www.googleapis.com/books/v1/volumes/" + id + "?key=" + apiKey

	clientTimeout := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return obj, err
	}
	req.Header.Set("User-Agent", "book-request-api-call")
	res, err := clientTimeout.Do(req)
	if err != nil {
		return obj, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}
