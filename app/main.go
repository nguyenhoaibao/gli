package app

import (
	"fmt"
	"log"
	"os"

	"github.com/nguyenhoaibao/gli/crawler"
	"github.com/nguyenhoaibao/gli/shell"
)

func Start() {
	sites, err := Sites()
	if err != nil {
		log.Fatal(err)
	}

	sh := shell.New(os.Stdout)

	for _, s := range sites {
		for _, c := range s.Categories {
			itemName := fmt.Sprintf("%s_item", s.Name)
			itemCrawler := crawler.NewItemCrawler(itemName, s.Item.UrlPattern)

			cName := fmt.Sprintf("%s_%s", s.Name, c.Name)
			cCrawler := crawler.NewCategoryCrawler(cName, c.Url, c.Limit, c.CachedInSeconds)

			handler := crawler.HandlerFunc(cCrawler, itemCrawler)
			sh.Register(cName, shell.HandlerFunc(handler))
		}
	}

	// generic handler
	// handler for all input command
	sh.RegisterGeneric(crawler.GenericHandlerFunc())

	if err := sh.Start(); err != nil {
		log.Fatal(err)
	}
}
