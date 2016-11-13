package crawler

import (
	"errors"
	"io"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func GetDocumentFromReader(r io.Reader) (*goquery.Document, error) {
	if r == nil {
		return nil, errors.New("Params is required")
	}

	return goquery.NewDocumentFromReader(r)
}

func HandlerFunc(itemsCrawler *itemsCrawler, itemCrawler *itemCrawler) func(args ...string) (io.Reader, error) {
	return func(args ...string) (io.Reader, error) {
		items, err := itemsCrawler.Crawl()
		if err != nil {
			return nil, err
		}

		if len(args) == 0 {
			return items.Render(), nil
		}

		i, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, err
		}

		id := items.ItemN(i)
		if id == "" {
			return nil, nil
		}

		item, err := itemCrawler.Crawl(id)
		if err != nil {
			return nil, err
		}
		return item.Render(), nil
	}
}
