package crawler

import (
	"errors"
	"io"

	"github.com/PuerkitoBio/goquery"
)

func GetDocumentFromReader(r io.Reader) (*goquery.Document, error) {
	if r == nil {
		return nil, errors.New("Params is required")
	}

	return goquery.NewDocumentFromReader(r)
}
