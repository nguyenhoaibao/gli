package tr_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nguyenhoaibao/gli/tr"
	"github.com/russross/blackfriday"
)

func TestRender(t *testing.T) {
	content, err := ioutil.ReadFile("../testdata/markdown_content.md")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(content))

	// out := blackfriday.MarkdownCommon(content)
	// for i := range out {
	// 	fmt.Print(string(out[i]))
	// }

	r := tr.TerminalRenderer(0)
	_ = blackfriday.Markdown(content, r, 0)

	// err = ioutil.WriteFile("../testdata/m.md", content, 0666)

	// fmt.Printf("%+q", c)
	// fmt.Printf("%+q", string(c))
}
