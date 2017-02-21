package yelp

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var allCategoriesConfig []CategoryConfig

func init() {
	file, err := ioutil.ReadFile("yelp/categories.json")

	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}

	json.Unmarshal(file, &allCategoriesConfig)
}
