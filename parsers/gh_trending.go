package parsers

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/nguyenhoaibao/gli/crawler"
)

type (
	ghTrendingParser struct{}

	ghTrendingItem struct {
		Title       string
		Owner       string
		Description string
		Url         string
		Language    string
		StarsText   string
	}

	ghTrendingItems []*ghTrendingItem
)

func NewGHTrendingParser() ghTrendingParser {
	var p ghTrendingParser
	return p
}

func init() {
	p := NewGHTrendingParser()
	crawler.Register("github_trending", p)
}

func (c ghTrendingParser) Parse(r io.Reader, limit int) (crawler.Renderer, error) {
	doc, err := crawler.GetDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	items := c.ParseItems(doc, limit)

	return items, nil
}

func (c ghTrendingParser) ParseItems(doc *goquery.Document, limit int) ghTrendingItems {
	var items ghTrendingItems

	qContainer := doc.Find(".repo-list")
	qitems := qContainer.Find(".repo-list-item")
	qitems.Each(func(_ int, s *goquery.Selection) {
		shortUrl, _ := s.Find(".repo-list-name a").Attr("href")
		title := string(shortUrl[1:])
		owner := s.Find("span.prefix").Text()
		description := strings.TrimSpace(s.Find("p.repo-list-description").Text())
		url := "https://github.com/" + title

		meta := strings.Split(s.Find(".repo-list-meta").Text(), "•")
		language := ""
		starsText := ""
		if l := len(meta); l == 3 {
			language = strings.TrimSpace(meta[0])
			starsText = strings.TrimSpace(meta[1])
		} else if l == 2 {
			starsText = strings.TrimSpace(meta[0])
		}

		items = append(items, &ghTrendingItem{
			Title:       title,
			Owner:       owner,
			Description: description,
			Url:         url,
			Language:    language,
			StarsText:   starsText,
		})
	})

	return items
}

func (items ghTrendingItems) Render() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n")

	for index, item := range items {
		si := fmt.Sprintf("%d", index+1)

		buffer.WriteString(fmt.Sprintf("%s.  %s", color.CyanString(si), color.YellowString(item.Title)))
		buffer.WriteString("\n")
		buffer.WriteString(fmt.Sprintf("    %s", color.MagentaString(item.Description)))
		buffer.WriteString("\n")
		if item.Language != "" {
			buffer.WriteString(fmt.Sprintf("    %s | %s", color.GreenString(item.Language), color.RedString(item.StarsText)))
		} else {
			buffer.WriteString(fmt.Sprintf("    %s", color.RedString(item.StarsText)))
		}

		buffer.WriteString("\n\n")
	}

	// fmt.Println(buffer.String())

	return buffer.String()
}
