package yelp

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type CategoryConfig struct {
	Alias            string   `json:"alias"`
	Title            string   `json:"title"`
	CountryWhitelist []string `json:"country_whitelist"`
	Parents          []string `json:"parents"`
}

type Category struct {
	Id         string
	Alias      string
	Title      string
	CreatedAt  string
	ModifiedAt string
}

type CategoriesConfig []CategoryConfig
type Categories []Category

var allCategories CategoriesConfig

func init() {
	file, err := ioutil.ReadFile("yelp/categories.json")

	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}

	json.Unmarshal(file, &allCategories)
}
