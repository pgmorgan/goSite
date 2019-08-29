package bookapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Volume struct {
	Title  string   `json:"title"`
	Author []string `json:"authors"`
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

func FindTopTen(title string) (JsonObjects, error) {
	uri := "https://www.googleapis.com/books/v1/volumes?q=" +
		url.QueryEscape(title) +
		"&key=AIzaSyCysdPDapS3TvjGhE1ZxMplLol4MQgR9Ks"
	obj := JsonObjects{}

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

func FindOne(id string) (JsonObject, error) {
	uri := "https://www.googleapis.com/books/v1/volumes/" + id +
		"?key=AIzaSyCysdPDapS3TvjGhE1ZxMplLol4MQgR9Ks"
	obj := JsonObject{}

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
