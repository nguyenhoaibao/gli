package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nguyenhoaibao/gli/crawler"
	_ "github.com/nguyenhoaibao/gli/parsers"
)

func main() {
	site := flag.String("s", "", "Sites to view")
	all := flag.Bool("all", false, "Run all configure sites")

	flag.Parse()

	if *all {
		crawler.RunAll()
		return
	}

	if *site != "" {
		crawler.Run(*site)
		return
	}

	fmt.Println("Not enough arguments. Please run -h for more information.")
	os.Exit(1)
}
