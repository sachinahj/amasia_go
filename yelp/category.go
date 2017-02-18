package yelp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)


type Category struct {
	Alias            string    `json:"alias"`
	Title            string    `json:"title"`
  CountryWhitelist []string  `json:"country_whitelist"`
	Parents []string  `json:"parents"`
}

type Categories []Category
var allCategories Categories

func init() {
  fmt.Println("initing yelp category")
	file, err := ioutil.ReadFile("yelp/categories.json")

  if err != nil {                // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error parsing yelp/categories.json: %s \n", err))
	}

  json.Unmarshal(file, &allCategories)
}
