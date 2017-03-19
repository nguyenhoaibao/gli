package parsers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/nguyenhoaibao/gli/cli"
	"github.com/nguyenhoaibao/gli/crawler"
	"github.com/russross/blackfriday"
)

const (
	GithubItemName = "github_item"

	commonTerminalFlags = 0 |
		cli.THeaderMarkdown |
		cli.THeaderAppendLine |
		cli.TCodeMarkdown
)

func init() {
	var p ghItemParser
	crawler.RegisterItemParser(GithubItemName, p)
}

//---------------------------------------------------

type (
	ghItem struct {
		Content string `json:"content"`
	}

	ghItemParser struct{}
)

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

	// fmt.Println(string(dec))

	// dec, _ := ioutil.ReadFile("./testdata/markdown_content.md")

	return ghItemContent(dec), nil
}

//---------------------------------------------------

type ghItemContent []byte

func (c ghItemContent) Render() io.Reader {
	// var buffer bytes.Buffer
	//
	// cmd := exec.Command("mdr")
	// cmd.Stdin = bytes.NewBuffer(b)
	// cmd.Stdout = &buffer
	//
	// // mdv command does not found
	// if err := cmd.Run(); err == nil {
	// 	return &buffer
	// }
	// return bytes.NewBuffer(b)

	t := cli.TerminalRenderer(commonTerminalFlags)
	content := blackfriday.Markdown(c, t, 0)

	return bytes.NewBuffer(content)
}
