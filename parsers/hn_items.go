package parsers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/olekukonko/tablewriter"
)

type hnItems struct {
	items []int
}

func NewHNItems(items []int, limit int) *hnItems {
	if len(items) > limit {
		items = items[:limit]
	}

	return &hnItems{items: items}
}

func ParseHNItemsFromReader(r io.Reader) ([]int, error) {
	var items []int
	err := json.NewDecoder(r).Decode(&items)

	return items, err
}

func (hn hnItems) Render() string {
	items := hn.items
	results := make(chan *hnItem)

	var wg sync.WaitGroup
	wg.Add(len(items))

	for _, id := range items {
		go func(id int) {
			defer wg.Done()

			result, err := ParseHNItem(id)
			if err != nil {
				log.Println(err)
				return
			}

			results <- result
		}(id)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var buffer bytes.Buffer
	var index int

	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{"#", "Title", "Points"})
	table.SetColWidth(90)

	for result := range results {
		index += 1
		table.Append([]string{
			fmt.Sprintf("%d", index),
			result.Title,
			fmt.Sprintf("%d", result.Score),
		})
	}

	table.Render()

	return ""

	// return &buffer
}
