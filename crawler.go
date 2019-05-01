package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getCrawlerURL() string {
	addr := os.Getenv("BOOK_PORT_8080_TCP_ADDR")
	port := os.Getenv("BOOK_PORT_8080_TCP_PORT")
	return fmt.Sprintf("http://%s:%s", addr, port)
}

func searchBook(keyword string) []*Book {
	ret := make([]*Book, 0)
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
	for index, book := range ret {
		ret[index] = bookManager.checkBook(book)
	}
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
