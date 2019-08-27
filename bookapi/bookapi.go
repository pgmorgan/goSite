package bookapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type volume struct {
	Title string `json:"title"`
}

type itemInfo struct {
	VolumeInfo volume `json:volumeInfo"`
	// Kind string `json:"kind"`
}

type jsonObject struct {
	Items []itemInfo `json:"items"`
}

func FindTopTen(title string) (jsonObject, error) {
	url := "https://www.googleapis.com/books/v1/volumes?q=harry+potter"
	obj := jsonObject{}

	clientTimeout := http.Client{
		Timeout: time.Second * 4,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
