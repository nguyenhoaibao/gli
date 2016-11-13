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

	s := shell.New(os.Stdout)

	for _, site := range sites {
		for _, items := range site.Types {
			itemName := fmt.Sprintf("%s_item", site.Name)
			itemCrawler := crawler.NewItemCrawler(itemName, site.Item.UrlPattern)

			itemsName := fmt.Sprintf("%s_%s", site.Name, items.Name)
			itemsCrawler := crawler.NewItemsCrawler(itemsName, items.Url, items.Limit, items.CachedInSeconds)

			handler := crawler.HandlerFunc(itemsCrawler, itemCrawler)
			s.Register(itemsName, shell.HandlerFunc(handler))
		}
	}

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
