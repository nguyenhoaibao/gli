package crawler

import (
	"fmt"
	"io"
	"log"

	"github.com/nguyenhoaibao/gli/app"
)

type Renderer interface {
	Render() string
}

type Crawler interface {
	Crawl() (Renderer, error)
}

type Parser interface {
	Parse(io.Reader, int) (Renderer, error)
}

var parsers = make(map[string]Parser)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Run(site *app.Site) (string, error) {
	items, err := Crawl(site)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return items.Render(), nil
}

func Crawl(site *app.Site) (Renderer, error) {
	fmt.Printf("Gathering data for site %s...", site.Name)
	fmt.Println()

	var crawler Crawler

	switch site.Type {
	case "html":
		crawler = NewHtmlCrawler(site)
	case "json":
		crawler = NewJsonCrawler(site)
	}

	return crawler.Crawl()
}

func Register(name string, parser Parser) {
	if _, exists := parsers[name]; exists {
		log.Fatalf("%s already exists", name)
	}
	parsers[name] = parser
}
