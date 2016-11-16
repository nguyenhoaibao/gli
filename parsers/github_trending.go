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

const GITHUB_TRENDING_NAME = "github_trending"

func init() {
	var p ghTrendingParser
	crawler.RegisterCategoryParser(GITHUB_TRENDING_NAME, p)
}

//---------------------------------------------------

type (
	ghTrendingResult struct {
		Title       string
		Owner       string
		Description string
		Url         string
		Language    string
		TotalStars  string
		TodayStars  string
	}

	ghTrendingResults []*ghTrendingResult
)

func (results ghTrendingResults) Total() int {
	return len(results)
}

func (results ghTrendingResults) ItemN(i int) string {
	if results.Total() < i {
		return ""
	}
	return results[i-1].Title
}

func (results ghTrendingResults) Render() io.Reader {
	var buffer bytes.Buffer
	total := len(results)

	buffer.WriteString("\n")

	for i, r := range results {
		var indent int

		if i < 9 {
			indent = 2
		} else {
			indent = 1
		}

		idxStr := fmt.Sprintf("%d", i+1)

		// print title
		fmt.Fprintf(&buffer, "%s.%s%s\n", color.CyanString(idxStr), strings.Repeat(" ", indent), color.YellowString(r.Title))

		// print description
		fmt.Fprintf(&buffer, "%s%s\n", strings.Repeat(" ", 4), color.MagentaString(r.Description))

		// print meta
		if r.Language != "" {
			fmt.Fprintf(&buffer, "%s%s | %s", strings.Repeat(" ", 4), color.GreenString(r.Language), color.RedString(r.TodayStars))
		} else {
			fmt.Fprintf(&buffer, "%s%s", strings.Repeat(" ", 4), color.RedString(r.TodayStars))
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

func (p ghTrendingParser) Parse(r io.Reader, limit int) (crawler.CategoryRenderer, error) {
	doc, err := crawler.GetDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return p.parse(doc, limit), nil
}

func (p ghTrendingParser) parse(doc *goquery.Document, limit int) ghTrendingResults {
	var results ghTrendingResults

	qContainer := doc.Find(".repo-list")
	qItems := qContainer.Find("li")
	qItems.Each(func(_ int, s *goquery.Selection) {
		if len(results) >= limit {
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

		results = append(results, &ghTrendingResult{
			Title:       title,
			Owner:       owner,
			Description: description,
			Url:         url,
			Language:    language,
			TotalStars:  totalStars,
			TodayStars:  todayStars,
		})
	})

	return results
}
