package cli_test

import (
	"io/ioutil"
	"testing"

	"github.com/nguyenhoaibao/gli/cli"
	"github.com/russross/blackfriday"
)

func TestRender(t *testing.T) {
	content, err := ioutil.ReadFile("../testdata/markdown_content.md")
	if err != nil {
		t.Fatal(err)
	}

	// out := blackfriday.MarkdownCommon(content)
	// for i := range out {
	// 	fmt.Print(string(out[i]))
	// }

	r := cli.TerminalRenderer(0)
	_ = blackfriday.Markdown(content, r, 0)

	// fmt.Println(string(c))
}
