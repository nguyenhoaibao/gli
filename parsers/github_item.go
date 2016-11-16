package parsers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"os/exec"

	"github.com/nguyenhoaibao/gli/crawler"
)

const GITHUB_ITEM_NAME = "github_item"

func init() {
	var p ghItemParser
	crawler.RegisterItemParser(GITHUB_ITEM_NAME, p)
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

//---------------------------------------------------

type ghItemContent []byte

func (b ghItemContent) Render() io.Reader {
	var buffer bytes.Buffer

	cmd := exec.Command("mdv", "-")
	cmd.Stdin = bytes.NewBuffer(b)
	cmd.Stdout = &buffer

	// mdv command does not found
	if err := cmd.Run(); err == nil {
		return &buffer
	}
	return bytes.NewBuffer(b)
}
