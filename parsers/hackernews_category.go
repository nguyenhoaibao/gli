package parsers

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

func parseHNCategory(doc *goquery.Document, limit int) hnCategoryItems {
	var items hnCategoryItems

	qContainer := doc.Find("#hnmain")
	qItems := qContainer.Find("table.itemlist tr")
	qItems.Each(func(i int, s *goquery.Selection) {
		total := len(items)

		if id, exists := s.Attr("id"); exists {
			link := s.Find("td.title > a")
			title := link.Text()
			url, _ := link.Attr("href")
			fromSite := link.Siblings().Find("a").Text()

			items = append(items, &hnCategoryItem{
				Id:       id,
				Title:    title,
				FromSite: fromSite,
				Url:      url,
			})
		} else if subtext := s.Find("td.subtext"); subtext.Length() > 0 {
			score := subtext.Find("span.score").Text()
			by := subtext.Find("a.hnuser").Text()

			var time, totalComments string

			qLinks := subtext.Find("a[href^=item]")
			qLinks.Each(func(_ int, link *goquery.Selection) {
				if link.ParentFiltered(".age").Length() > 0 {
					time = link.Text()
				} else if link.ParentFiltered(".subtext").Length() > 0 {
					totalComments = link.Text()
				}
			})

			items[total-1].Score = score
			items[total-1].By = by
			items[total-1].Time = time
			items[total-1].TotalComments = totalComments
		}
	})

	items = items[:limit]
	return items
}

//---------------------------------------------------

type (
	hnCategoryItem struct {
		Id            string
		Title         string
		FromSite      string
		Url           string
		Score         string
		By            string
		Time          string
		TotalComments string
	}

	hnCategoryItems []*hnCategoryItem
)

func (items hnCategoryItems) Total() int {
	return len(items)
}

func (items hnCategoryItems) ItemN(i int) string {
	if i > items.Total() {
		return ""
	}
	return items[i-1].Id
}

func (items hnCategoryItems) Render() io.Reader {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "\n")

	for i, item := range items {
		var indent int

		if i < 9 {
			indent = 2
		} else {
			indent = 1
		}

		iStr := fmt.Sprintf("%d", i+1)

		// print title
		fmt.Fprintf(&buffer, "%s.%s%s", color.CyanString(iStr), strings.Repeat(" ", indent), color.YellowString(item.Title))
		if item.FromSite != "" {
			fmt.Fprintf(&buffer, " (%s)", color.GreenString(item.FromSite))
		}
		fmt.Fprintf(&buffer, "\n")

		// print meta info
		if item.Score != "" && item.By != "" {
			fmt.Fprintf(&buffer, "%s%s by %s %s", strings.Repeat(" ", 4), color.RedString(item.Score), color.BlueString(item.By), color.MagentaString(item.Time))
		} else {
			fmt.Fprintf(&buffer, "%s%s", strings.Repeat(" ", 4), color.MagentaString(item.Time))
		}
		if item.TotalComments != "" {
			fmt.Fprintf(&buffer, " | %s", item.TotalComments)
		}
		fmt.Fprintf(&buffer, "\n")

		if i != items.Total()-1 {
			fmt.Fprint(&buffer, "\n")
		}
	}

	return &buffer
}
