package tr

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/russross/blackfriday"
)

const (
	ListItemPadding = 2

	BlockCodePadding = 4
	BlockCodeFormat  = "```"

	InlineCodeFormat = "`"
)

const (
	THeaderWithLevel = 1 << iota // Header with level
	THeaderWithLineSuffix

	TCodeWithBacktick // Code with backtick prefix and suffix
)

var (
	listIndex = 1

	Style = map[string][]color.Attribute{
		"Heading":     []color.Attribute{color.FgYellow, color.Bold},
		"Link":        []color.Attribute{color.FgCyan},
		"LinkContent": []color.Attribute{color.FgBlue, color.Bold, color.Underline},
		"CodeSpan":    []color.Attribute{color.FgHiMagenta},
	}
)

type Terminal struct {
	flags int // T* options
}

func TerminalRenderer(flags int) blackfriday.Renderer {
	return &Terminal{
		flags: flags,
	}
}

func repeat(out *bytes.Buffer, s string, n int) {
	out.WriteString(fmt.Sprintf("%s", strings.Repeat(s, n)))
}

func drawline(out *bytes.Buffer) {
	repeat(out, "-", 50)
}

func newline(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}

func doubleNewline(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteString("\n\n")
	}
}

func tab(out *bytes.Buffer, padding int) {
	repeat(out, " ", padding)
}

func (t *Terminal) styleInlineCode(out *bytes.Buffer) {
	if t.flags&TCodeWithBacktick == 0 {
		return
	}
	out.WriteString(InlineCodeFormat)
}

func (t *Terminal) styleBlockCode(out *bytes.Buffer) {
	if t.flags&TCodeWithBacktick == 0 {
		return
	}
	tab(out, BlockCodePadding)
	out.WriteString(BlockCodeFormat)
}

func (t *Terminal) writeInlineCode(out *bytes.Buffer, text string) {
	t.styleInlineCode(out)
	defer t.styleInlineCode(out)

	out.WriteString(text)
}

func (t *Terminal) writeBlockCode(out *bytes.Buffer, text string) {
	t.styleBlockCode(out)

	defer func() {
		t.styleBlockCode(out)
		newline(out)
	}()

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		tab(out, BlockCodePadding)
		out.WriteString(string(line))
		newline(out)
	}
}

// Block-level callbacks
func (t *Terminal) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	fmt.Println("BlockCode")
	out.WriteString("Block code here")
}
func (t *Terminal) BlockQuote(out *bytes.Buffer, text []byte) {
	fmt.Println("BlockQuote")
}
func (t *Terminal) BlockHtml(out *bytes.Buffer, text []byte) {
	fmt.Println("BlockHtml")
}

func (t *Terminal) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	marker := out.Len()
	doubleNewline(out)

	color.Output = out
	color.Set(Style["Heading"]...)
	defer color.Unset()

	if t.flags&THeaderWithLevel != 0 {
		repeat(out, "#", level)
		out.WriteByte(' ')
	}

	// fmt.Println(out.Bytes())

	if !text() {
		out.Truncate(marker)
		return
	}

	// fmt.Println(out.Bytes())
	// fmt.Printf("%+q", string(out.Bytes()))

	if level == 1 || level == 2 {
		newline(out)
		if t.flags&THeaderWithLineSuffix != 0 {
			drawline(out)
		}
	}
}

func (t *Terminal) HRule(out *bytes.Buffer) {
	fmt.Println("HRule")
}

func (t *Terminal) ListItem(out *bytes.Buffer, text []byte, flags int) {
	if len(text) == 0 {
		return
	}
	tab(out, ListItemPadding)

	if flags&blackfriday.LIST_TYPE_ORDERED == blackfriday.LIST_TYPE_ORDERED {
		out.WriteString(fmt.Sprintf("%d. ", listIndex))
		listIndex++
	} else {
		out.WriteString("* ")
	}
	out.Write(text)
	newline(out)
}

func (t *Terminal) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	// newline(out)
	doubleNewline(out)

	if !text() {
		out.Truncate(marker)
		return
	}
}
func (t *Terminal) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	fmt.Println("Table")
}
func (t *Terminal) TableRow(out *bytes.Buffer, text []byte) {
	fmt.Println("TableRow")
}
func (t *Terminal) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
	fmt.Println("TableHeaderCell")
}
func (t *Terminal) TableCell(out *bytes.Buffer, text []byte, flags int) {
	fmt.Println("TableCell")
}
func (t *Terminal) Footnotes(out *bytes.Buffer, text func() bool) {
	fmt.Println("Footnotes")
}
func (t *Terminal) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	fmt.Println("FootnoteItem")
}
func (t *Terminal) TitleBlock(out *bytes.Buffer, text []byte) {
	fmt.Println("TitleBlock")
}

// Span-level callbacks
func (t *Terminal) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	fmt.Println("AutoLink")
}

func (t *Terminal) CodeSpan(out *bytes.Buffer, text []byte) {
	color.Output = out
	color.Set(Style["CodeSpan"]...)
	defer color.Unset()

	s := string(text)
	if strings.Index(s, "\n") < 0 {
		t.writeInlineCode(out, s)
		return
	}
	t.writeBlockCode(out, s)
}

func (t *Terminal) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	fmt.Println("DoubleEmphasis")
}

func (t *Terminal) Emphasis(out *bytes.Buffer, text []byte) {
	fmt.Println("Emphasis", string(text))
}

func (t *Terminal) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	var name []byte

	if len(alt) > 0 {
		name = alt
	} else if len(title) > 0 {
		name = title
	} else {
		name = []byte("Image")
	}

	out.Write(name)
}

func (t *Terminal) LineBreak(out *bytes.Buffer) {
	fmt.Println("LineBreak")
}

func (t *Terminal) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	color.Output = out

	if len(content) == 0 {
		color.Set(Style["Link"]...)
		out.Write(link)
		color.Unset()
		return
	}

	color.Set(Style["LinkContent"]...)
	out.Write(content)
	color.Unset()

	color.Set(Style["Link"]...)
	out.WriteString(fmt.Sprintf("(%s)", string(link)))
	color.Unset()
}

func (t *Terminal) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	newline(out)

	listIndex = 1

	if !text() {
		out.Truncate(marker)
		return
	}
}

func (t *Terminal) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	fmt.Println("RawHtmlTag")
}

func (t *Terminal) TripleEmphasis(out *bytes.Buffer, text []byte) {
	fmt.Println("TripleEmphasis")
}

func (t *Terminal) StrikeThrough(out *bytes.Buffer, text []byte) {
	fmt.Println("StrikeThrough")
}

func (t *Terminal) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	fmt.Println("FootnoteRef")
}

// Low-level callbacks
func (t *Terminal) Entity(out *bytes.Buffer, entity []byte) {
	fmt.Println("Entity")
}

func (t *Terminal) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

// Header and footer
func (t *Terminal) DocumentHeader(out *bytes.Buffer) {
	fmt.Println("DocumentHeader")
}
func (t *Terminal) DocumentFooter(out *bytes.Buffer) {
	fmt.Println("DocumentFooter")
}

func (t *Terminal) GetFlags() int {
	return 0
}
