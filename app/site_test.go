package app_test

import (
	"testing"

	"github.com/nguyenhoaibao/gotools/app"
)

func TestGetSites(t *testing.T) {
	sites, err := app.GetSites()
	if err != nil {
		t.Fatal(err)
	}

	if len(sites) <= 0 {
		t.Fatal("Total sites should greater than 0")
	}
}
