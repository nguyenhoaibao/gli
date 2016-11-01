package crawler

import (
	"errors"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Query(url string) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("Url is required")
	}

	return http.Get(url)
}

func GetDocumentFromReader(r io.Reader) (*goquery.Document, error) {
	if r == nil {
		return nil, errors.New("Params is required")
	}

	return goquery.NewDocumentFromReader(r)
}
