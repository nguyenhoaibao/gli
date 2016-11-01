package parsers

import (
	"fmt"
	"io"

	"github.com/nguyenhoaibao/gli/crawler"
)

const hnNewName = "hackernews_new"

type hnNewParser struct{}

func init() {
	var p hnNewParser
	crawler.Register(hnNewName, p)
}

func (c hnNewParser) Parse(r io.Reader, limit int) (crawler.Renderer, error) {
	items, err := ParseHNItemsFromReader(r)
	if err != nil {
		return nil, err
	}
	if len(items) <= 0 {
		return nil, fmt.Errorf("Cannot parse any items for %s", hnNewName)
	}

	hnItems := NewHNItems(items, limit)

	return hnItems, nil
}
