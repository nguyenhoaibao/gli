package app

import (
	"encoding/json"
	"fmt"
	"os"
)

const dataFile = "data/sites.json"

type Site struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	Url       string `json:"url"`
}

func GetSites() ([]*Site, error) {
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

func GetSiteByName(siteName string) (*Site, error) {
	sites, err := GetSites()
	if err != nil {
		return nil, err
	}

	for _, site := range sites {
		if siteName == site.Name || siteName == site.ShortName {
			return site, nil
		}
	}

	return nil, fmt.Errorf("Cannot get any site with name %s", siteName)
}
