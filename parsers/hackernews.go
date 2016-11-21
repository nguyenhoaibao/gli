package parsers

import (
	"io"

	"github.com/nguyenhoaibao/gli/crawler"
)

type hnParser struct{}

var hnCategories = []string{
	"hackernews_top",
	"hackernews_new",
	"hackernews_show",
	"hackernews_ask",
	"hackernews_jobs",
}

func init() {
	for _, c := range hnCategories {
		var p hnParser
		crawler.RegisterCategoryParser(c, p)
	}
}

func (c hnParser) Parse(r io.Reader, limit int) (crawler.CategoryRenderer, error) {
	doc, err := crawler.DocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return parseHNCategory(doc, limit), nil
}
