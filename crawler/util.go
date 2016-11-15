package crawler

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetDocumentFromReader(r io.Reader) (*goquery.Document, error) {
	if r == nil {
		return nil, errors.New("Params is required")
	}

	return goquery.NewDocumentFromReader(r)
}

func HandlerFunc(categoryCrawler *categoryCrawler, itemCrawler *itemCrawler) func(args ...string) (io.Reader, error) {
	return func(args ...string) (io.Reader, error) {
		fmt.Print("Please wait")

		ch := make(chan io.Reader)
		chErr := make(chan error)

		go func() {
			items, err := categoryCrawler.Crawl()
			if err != nil {
				chErr <- err
				return
			}

			if len(args) == 0 {
				ch <- items.Render()
				return
			}

			i, err := strconv.Atoi(args[0])
			if err != nil {
				chErr <- err
				return
			}

			id := items.ItemN(i)
			if id == "" {
				ch <- nil
				return
			}

			item, err := itemCrawler.Crawl(id)
			if err != nil {
				chErr <- err
				return
			}
			ch <- item.Render()
		}()

		for {
			select {
			case result := <-ch:
				return result, nil
			case err := <-chErr:
				return nil, err
			case <-time.After(500 * time.Millisecond):
				fmt.Print(".")
			}
		}
	}
}
