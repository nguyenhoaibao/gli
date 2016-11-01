package parsers

import (
	"fmt"
	"io"

	"github.com/nguyenhoaibao/gli/crawler"
)

const hnJobName = "hackernews_job"

type hnJobParser struct{}

func init() {
	var p hnJobParser
	crawler.Register(hnJobName, p)
}

func (c hnJobParser) Parse(r io.Reader, limit int) (crawler.Renderer, error) {
	items, err := ParseHNItemsFromReader(r)
	if err != nil {
		return nil, err
	}
	if len(items) <= 0 {
		return nil, fmt.Errorf("Cannot parse any items for %s", hnJobName)
	}

	hnItems := NewHNItems(items, limit)

	return hnItems, nil
}
