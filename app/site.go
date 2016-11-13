package app

import (
	"encoding/json"
	"os"
	"time"
)

const dataFile = "data/sites.json"

type Items struct {
	Name            string        `json:"name"`
	Url             string        `json:"url"`
	Limit           int           `json:"limit"`
	CachedInSeconds time.Duration `json:"cached_in_seconds"`
}

type Item struct {
	UrlPattern string `json:"url_pattern"`
}

type Site struct {
	Name  string   `json:"name"`
	Types []*Items `json:"types"`
	Item  Item     `json:"item"`
}

func Sites() ([]*Site, error) {
	f, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	var sites []*Site
	err = json.NewDecoder(f).Decode(&sites)
	if err != nil {
		return nil, err
	}

	return sites, nil
}
