package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getCrawlerURL() string {
	return "http://127.0.0.1:8080"
}

func searchBook(keyword string) []Book {
	ret := make([]Book, 0)
	url := fmt.Sprintf("%s/search/%s", getCrawlerURL(), keyword)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return ret
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return ret
	}
	json.Unmarshal(body, &ret)
	return ret
}

func fetchBook(platform string, bookID int64) []Chapter {
	ret := make([]Chapter, 0)
	url := fmt.Sprintf("%s/fetch/%s/%d", getCrawlerURL(), platform, bookID)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return ret
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return ret
	}
	json.Unmarshal(body, &ret)
	return ret
}

func markBook(platform string, bookID int64) {
	url := fmt.Sprintf("%s/mark/%s/%d", getCrawlerURL(), platform, bookID)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return
	}
	resp.Body.Close()
	return
}
