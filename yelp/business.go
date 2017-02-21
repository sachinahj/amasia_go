package yelp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/viper"
)

func BusinessesSearch(zc ZipCode, lg LogConfigBusinessesSearch) {
	viper.SetConfigName("_config") // name of config file (without extension)
	viper.AddConfigPath(".")       // optionally look for config in the working directory
	err := viper.ReadInConfig()    // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		log.Fatal(err)
	}

	baseUrl := viper.Get("dataProviders.yelp.baseUrl")
	rawUrl := fmt.Sprintf("%v/%v", baseUrl, "businesses/search")
	accessToken := viper.Get("dataProviders.yelp.accessToken")

	req, err := http.NewRequest("GET", rawUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	q := req.URL.Query()
	q.Add("location", strconv.Itoa(zc.ZipCode))
	q.Add("sort_by", "rating")
	q.Add("offset", strconv.Itoa(lg.Offset))
	q.Add("limit", strconv.Itoa(lg.Limit))
	q.Add("categories", lg.Alias)
	req.URL.RawQuery = q.Encode()

	fmt.Println("req.URL.String()", req.URL.String())

	client := &http.Client{}

	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var bsr = new(BusinessesSearchResponse)
	err = json.Unmarshal([]byte(body), &bsr)
	if err != nil {
		log.Fatal(err)
	}

	var categories = make(map[string]Category)
	var businessCategories []BusinessCategory

	for _, b := range bsr.Businesses {
		b.ZipCode = zc.ZipCode

		fmt.Println("---------------------------")
		fmt.Println(b)
		b.Update()

		for _, c := range b.Categories {
			bc := BusinessCategory{BusinessId: b.Id, CategoryAlias: c.Alias}
			businessCategories = append(businessCategories, bc)
			categories[c.Alias] = c
			fmt.Println(bc)
		}
	}

	fmt.Println("---------------------------")
	fmt.Println(categories)
	for _, c := range categories {
		fmt.Println(c)
		c.Update()
	}

	fmt.Println("---------------------------")
	for _, bc := range businessCategories {
		fmt.Println(bc)
		bc.Update()
	}
}
