package parsers

import (
	"bufio"
	"os"
	"testing"
)

func TestParseGithubItem(t *testing.T) {
	f, err := os.Open("../testdata/github_item.json")
	if err != nil {
		t.Error(err)
	}

	var (
		p ghItemParser
		r = bufio.NewReader(f)
	)

	item, err := p.Parse(r)
	if err != nil {
		t.Error(err)
	}
	if item == nil {
		t.Error("Cannot parse github item")
	}
	item.Render()
}
