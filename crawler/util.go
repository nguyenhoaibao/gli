package crawler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Query(url string) (*http.Response, error) {
	if url == "" {
		return nil, errors.New("Url is required")
	}

	ch := make(chan *http.Response)
	chErr := make(chan error)

	fmt.Printf("Querying url %s", url)

	go func() {
		resp, err := http.Get(url)
		if err != nil {
			chErr <- err
			return
		}
		ch <- resp
	}()

	for {
		select {
		case resp := <-ch:
			return resp, nil
		case err := <-chErr:
			return nil, err
		case <-time.After(500 * time.Millisecond):
			fmt.Print(".")
		}
	}
}

func GetDocumentFromReader(r io.Reader) (*goquery.Document, error) {
	if r == nil {
		return nil, errors.New("Params is required")
	}

	return goquery.NewDocumentFromReader(r)
}
