package app

import (
	"encoding/json"
	"os"
)

const dataFile = "data/sites.json"

type Site struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Limit int    `json:"limit"`
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
