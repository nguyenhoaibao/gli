package main

import (
	"io"
	"log"
	"os"

	"github.com/nguyenhoaibao/gli/app"
	"github.com/nguyenhoaibao/gli/crawler"
	_ "github.com/nguyenhoaibao/gli/parsers"
	"github.com/nguyenhoaibao/gli/shell"
)

func createHandlerFunc(site *app.Site) func(args ...string) (io.Reader, error) {
	return func(args ...string) (io.Reader, error) {
		c := crawler.New(site.Name, site.Url, site.Limit)
		result, err := c.Crawl()
		if err != nil {
			return nil, err
		}
		return result.Render(), nil
	}
}

func main() {
	sites, err := app.Sites()
	if err != nil {
		log.Fatal(err)
	}

	shell := shell.New(os.Stdout)

	for _, site := range sites {
		shell.Register(site.Name, createHandlerFunc(site))
	}

	if err := shell.Start(); err != nil {
		log.Fatal(err)
	}
}
