package crawler_test

import (
	"testing"

	"github.com/nguyenhoaibao/gli/crawler"
)

func TestQueryWithUrlEmpty(t *testing.T) {
	url := ""
	_, err := crawler.Query(url)
	if err == nil {
		t.Error("Should receive error when url is empty")
	}
}

func TestQueryWithUrlInvalid(t *testing.T) {
	url := "abcde://123"
	_, err := crawler.Query(url)
	if err == nil {
		t.Error("Should receive error when url is invalid")
	}
}

// func TestQueryWithUrlValid(t *testing.T) {
// 	url := "http://google.com"
// 	_, err := crawler.Query(url)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
