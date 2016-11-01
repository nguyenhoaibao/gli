package main

import (
	"log"
	"os"

	"github.com/abiosoft/ishell"
	"github.com/nguyenhoaibao/gli/app"
	"github.com/nguyenhoaibao/gli/crawler"
	_ "github.com/nguyenhoaibao/gli/parsers"
)

func main() {
	sites, err := app.Sites()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	shell := ishell.New()
	shell.Println("Welcome")

	for _, site := range sites {
		go func(s *app.Site) {
			shell.Register(s.Name, func(args ...string) (string, error) {
				return crawler.Run(s)
			})
		}(site)
	}

	shell.Start()
}
