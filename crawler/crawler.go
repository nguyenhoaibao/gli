package crawler

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/nguyenhoaibao/gli/app"
)

type Crawler interface {
	Crawl() (io.Writer, error)
}

type Parser interface {
	Parse(io.Reader) (io.Writer, error)
}

var parsers = make(map[string]Parser)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Run(siteName string) {
	site, err := app.GetSiteByName(siteName)
	if err != nil {
		log.Fatal(err)
	}

	results := make(chan io.Writer)

	go func() {
		Crawl(site, results)
		close(results)
	}()

	DisplayResults(results)
}

func RunAll() {
	sites, err := app.GetSites()
	if err != nil {
		log.Fatal(err)
	}
	if len(sites) == 0 {
		log.Fatal("Cannot load any sites")
	}

	results := make(chan io.Writer)

	var wg sync.WaitGroup
	wg.Add(len(sites))

	for _, site := range sites {
		go func(site *app.Site) {
			Crawl(site, results)
			wg.Done()
		}(site)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	DisplayResults(results)
}

func Crawl(site *app.Site, results chan<- io.Writer) {
	var crawler Crawler

	switch site.Type {
	case "html":
		crawler = NewHtmlCrawler(site)
		result, err := crawler.Crawl()
		if err != nil {
			log.Fatal(err)
		}

		results <- result
	}
}

func Register(name string, parser Parser) {
	if _, exists := parsers[name]; exists {
		log.Fatalf("%s already exists", name)
	}
	parsers[name] = parser
}

func DisplayResults(results <-chan io.Writer) {
	for result := range results {
		fmt.Println()
		fmt.Println(result)
	}
}
