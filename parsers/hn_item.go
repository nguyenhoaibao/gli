package parsers

import (
	"encoding/json"
	"fmt"

	"github.com/nguyenhoaibao/gli/crawler"
)

const ITEM_URL = "https://hacker-news.firebaseio.com/v0/item/%d.json"

type hnItem struct {
	Title string `json:"title"`
	Score int    `json:"score"`
	Time  int    `json:"time"`
	Url   string `json:"url"`
}

func ParseHNItemByUrl(url string) (*hnItem, error) {
	resp, err := crawler.Query(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var item hnItem
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func ParseHNItem(id int) (*hnItem, error) {
	url := fmt.Sprintf(ITEM_URL, id)

	return ParseHNItemByUrl(url)
}
