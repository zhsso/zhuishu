package main

import (
	"fmt"

	"../novel"
	"github.com/anaskhan96/soup"
)

// Search 搜索
func Search(keyword string) []novel.SimpleInfo {
	ret := make([]novel.SimpleInfo, 0)
	resp, err := soup.Get(fmt.Sprintf("https://www.qidian.com/search?kw=%s", keyword))
	if err != nil {
		return ret
	}
	doc := soup.HTMLParse(resp)
	books := doc.FindAll("div", "class", "book-mid-info")
	for _, book := range books {
		title := book.Find("a")
		println(title.Text())
		// println(title.Attrs()["data-bid"])
	}
	return ret
}

func main() {
	Search("修真聊天群")
}
