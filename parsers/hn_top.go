package parsers

import (
	"fmt"
	"io"

	"github.com/nguyenhoaibao/gli/crawler"
)

const hnTopName = "hackernews_top"

type hnTopParser struct{}

func init() {
	var p hnTopParser
	crawler.Register(hnTopName, p)
}

func (c hnTopParser) Parse(r io.Reader, limit int) (crawler.Renderer, error) {
	items, err := ParseHNItemsFromReader(r)
	if err != nil {
		return nil, err
	}
	if len(items) <= 0 {
		return nil, fmt.Errorf("Cannot parse any items for %s", hnTopName)
	}

	hnItems := NewHNItems(items, limit)

	return hnItems, nil
}
