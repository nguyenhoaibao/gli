package crawler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"time"
)

type Renderer interface {
	Render() *bytes.Buffer
}

type Parser interface {
	Parse(io.Reader, int) (Renderer, error)
}

type crawler struct {
	name  string
	url   string
	limit int
}

var parsers = make(map[string]Parser)

func New(name string, url string, limit int) *crawler {
	return &crawler{name, url, limit}
}

func Register(name string, parser Parser) error {
	if _, exists := parsers[name]; exists {
		return fmt.Errorf("%s already exists", name)
	}
	parsers[name] = parser
	return nil
}

func mockServer(content string) *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")
		w.WriteHeader(200)
		fmt.Fprintf(w, content)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func (c *crawler) Crawl() (Renderer, error) {
	fmt.Printf("Gathering data for site %s", c.name)

	chResult := make(chan Renderer)
	chErr := make(chan error)

	go func() {
		resp, err := c.Download()
		if err != nil {
			chErr <- err
			return
		}
		defer resp.Body.Close()

		result, err := c.Parse(resp.Body)
		if err != nil {
			chErr <- err
			return
		}
		chResult <- result
	}()

	for {
		select {
		case result := <-chResult:
			return result, nil
		case err := <-chErr:
			return nil, err
		case <-time.After(500 * time.Millisecond):
			fmt.Print(".")
		}
	}
}

func (c *crawler) Download() (*http.Response, error) {
	// return http.Get(c.url)
	content, err := ioutil.ReadFile(filepath.Join("testdata", c.name+".html"))
	if err != nil {
		return nil, err
	}

	server := mockServer(string(content[:]))
	defer server.Close()

	return http.Get(server.URL)
}

func (c *crawler) Parse(body io.Reader) (Renderer, error) {
	parser, exists := parsers[c.name]
	if !exists {
		return nil, fmt.Errorf("Parser %s does not exist", c.name)
	}
	return parser.Parse(body, c.limit)
}
