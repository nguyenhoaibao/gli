package crawler

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	currentCategoryCrawler *categoryCrawler
	currentItemCrawler     *itemCrawler
)

func setCurrentCrawler(cc *categoryCrawler, ic *itemCrawler) {
	currentCategoryCrawler = cc
	currentItemCrawler = ic
}

func getCurrentCrawler() (*categoryCrawler, *itemCrawler) {
	return currentCategoryCrawler, currentItemCrawler
}

func handlerFunc(cc *categoryCrawler, ic *itemCrawler, args ...string) (io.Reader, error) {
	fmt.Print("Please wait")

	ch := make(chan io.Reader)
	chErr := make(chan error)

	setCurrentCrawler(cc, ic)

	go func() {
		items, err := cc.Crawl()
		if err != nil {
			chErr <- err
			return
		}

		if len(args) == 0 {
			ch <- items.Render()
			return
		}

		if ic == nil {
			ch <- nil
			return
		}

		// get item at index i
		// if index is not int, return nil result instead of error
		i, err := strconv.Atoi(args[0])
		if err != nil {
			ch <- nil
			return
		}

		id := items.ItemN(i)
		if id == "" {
			ch <- nil
			return
		}

		item, err := ic.Crawl(id)
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

func DocumentFromReader(r io.Reader) (*goquery.Document, error) {
	if r == nil {
		return nil, errors.New("Params is required")
	}

	return goquery.NewDocumentFromReader(r)
}

func HandlerFunc(cc *categoryCrawler, ic *itemCrawler) func(args ...string) (io.Reader, error) {
	return func(args ...string) (io.Reader, error) {
		return handlerFunc(cc, ic, args...)
	}
}

func GenericHandlerFunc() func(args ...string) (io.Reader, error) {
	return func(args ...string) (io.Reader, error) {
		cc, ic := getCurrentCrawler()
		if cc == nil {
			return nil, nil
		}

		// check to see if first arg is number or not
		// if first arg is number, get data of current category for that index number
		// if not, return nil
		_, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, nil
		}

		return handlerFunc(cc, ic, args...)
	}
}
