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

func init() {
	var p ghTrendingParser
	crawler.RegisterItemsParser("github_trending", p)
}

//---------------------------------------------------

type (
	ghTrendingItem struct {
		Title       string
		Owner       string
		Description string
		Url         string
		Language    string
		TotalStars  string
		TodayStars  string
	}

	ghTrendingItems []*ghTrendingItem
)

func (items ghTrendingItems) Total() int {
	return len(items)
}

func (items ghTrendingItems) ItemN(i int) string {
	if items.Total() < i {
		return ""
	}
	return items[i-1].Title
}

func (items ghTrendingItems) Render() io.Reader {
	var buffer bytes.Buffer
	total := len(items)

	buffer.WriteString("\n")

	for i, item := range items {
		var indent int

		if i < 9 {
			indent = 2
		} else {
			indent = 1
		}

		idxStr := fmt.Sprintf("%d", i+1)

		// print title
		fmt.Fprintf(&buffer, "%s.%s%s\n", color.CyanString(idxStr), strings.Repeat(" ", indent), color.YellowString(item.Title))

		// print description
		fmt.Fprintf(&buffer, "%s%s\n", strings.Repeat(" ", 4), color.MagentaString(item.Description))

		// print meta
		if item.Language != "" {
			fmt.Fprintf(&buffer, "%s%s | %s", strings.Repeat(" ", 4), color.GreenString(item.Language), color.RedString(item.TodayStars))
		} else {
			fmt.Fprintf(&buffer, "%s%s", strings.Repeat(" ", 4), color.RedString(item.TodayStars))
		}

		if i == total-1 {
			fmt.Fprint(&buffer, "\n")
		} else {
			fmt.Fprint(&buffer, "\n\n")
		}
	}
	return &buffer
}

//---------------------------------------------------

type ghTrendingParser struct{}

func (p ghTrendingParser) Parse(r io.Reader, limit int) (crawler.ItemsRenderer, error) {
	doc, err := crawler.GetDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return p.parse(doc, limit), nil
}

func (p ghTrendingParser) parse(doc *goquery.Document, limit int) ghTrendingItems {
	var items ghTrendingItems

	qContainer := doc.Find(".repo-list")
	qItems := qContainer.Find("li")
	qItems.Each(func(_ int, s *goquery.Selection) {
		if len(items) >= limit {
			return
		}

		shortUrl, _ := s.Find("h3 > a").Attr("href")
		title := string(shortUrl[1:])
		owner := title[:strings.Index(title, "/")]
		description := strings.TrimSpace(s.Find("div.py-1 > p.d-inline-block").Text())
		url := "https://github.com/" + title

		qMeta := s.Find("div.mt-2")
		language := strings.TrimSpace(qMeta.Find("span[itemprop=programmingLanguage]").Text())
		totalStars := strings.TrimSpace(qMeta.Find("a[aria-label=Stargazers]").Text())
		todayStars := strings.TrimSpace(qMeta.Find("span.float-right").Text())

		items = append(items, &ghTrendingItem{
			Title:       title,
			Owner:       owner,
			Description: description,
			Url:         url,
			Language:    language,
			TotalStars:  totalStars,
			TodayStars:  todayStars,
		})
	})

	return items
}
