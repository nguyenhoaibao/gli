package parsers_test

import (
	"io/ioutil"
	"testing"

	"github.com/nguyenhoaibao/gli/parsers"
)

func TestParseHNItemWithIdInvalid(t *testing.T) {
	server := mockServer("null")
	defer server.Close()

	item, err := parsers.ParseHNItemByUrl(server.URL)
	if err != nil {
		t.Error(err)
	}
	if item.Title != "" {
		t.Error("Title should be empty when id is not valid")
	}
}

func TestParseHNItemWithIdValid(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/hn_item.json")
	if err != nil {
		t.Error(err)
	}

	server := mockServer(string(content[:]))
	defer server.Close()

	item, err := parsers.ParseHNItemByUrl(server.URL)
	if err != nil {
		t.Error(err)
	}
	if item.Title == "" {
		t.Error("Title should be not empty when id is valid")
	}
}
