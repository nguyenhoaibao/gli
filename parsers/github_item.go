package parsers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/nguyenhoaibao/gli/crawler"
)

func init() {
	var p ghItemParser
	crawler.RegisterItemParser("github_item", p)
}

//---------------------------------------------------

type ghItemContent []byte

func (b ghItemContent) Render() io.Reader {
	return bytes.NewBuffer(b)
}

//---------------------------------------------------

type ghItem struct {
	Content string `json:"content"`
}

//---------------------------------------------------

type ghItemParser struct{}

func (p ghItemParser) Parse(r io.Reader) (crawler.ItemRenderer, error) {
	var item ghItem
	err := json.NewDecoder(r).Decode(&item)
	if err != nil {
		return nil, err
	}

	dec, err := base64.StdEncoding.DecodeString(item.Content)
	if err != nil {
		return nil, err
	}

	return ghItemContent(dec), nil
}
