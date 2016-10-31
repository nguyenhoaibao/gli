package parsers

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/nguyenhoaibao/gli/crawler"
	"github.com/olekukonko/tablewriter"
)

type GithubTrendingResult struct {
	Title       string
	Owner       string
	Description string
	Url         string
	Language    string
	StarsText   string
}

type GithubTrendingParser struct{}

func init() {
	var parser GithubTrendingParser
	crawler.Register("github_trending", parser)
}

func (c GithubTrendingParser) Parse(r io.Reader) (io.Writer, error) {
	doc, err := crawler.GetDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	results := c.ParseResults(doc)

	return c.Display(results), nil
}

func (c GithubTrendingParser) ParseResults(doc *goquery.Document) []*GithubTrendingResult {
	var results []*GithubTrendingResult

	container := doc.Find(".repo-list")
	items := container.Find(".repo-list-item")
	items.Each(func(_ int, s *goquery.Selection) {
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

		results = append(results, &GithubTrendingResult{
			Title:       title,
			Owner:       owner,
			Description: description,
			Url:         url,
			Language:    language,
			StarsText:   starsText,
		})
	})

	return results
}

func (c GithubTrendingParser) Display(results []*GithubTrendingResult) io.Writer {
	var buffer bytes.Buffer

	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{"#", "Title", "Language", "Stars"})

	for i, result := range results {
		table.Append([]string{
			fmt.Sprintf("%d", i+1),
			result.Title,
			result.Language,
			result.StarsText,
		})
	}

	table.Render()

	return &buffer
}
