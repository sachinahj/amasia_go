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

type BusinessesSearchResponse struct {
	Total      int64      `json:"total"`
	Businesses []Business `json:"businesses"`
	Error				struct {
		Description string `json:"description"`
		Code string `json:"code"`
	} 		`json:"error"`
}

func BusinessesSearch(lg *Log) {
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

	lgc := lg.GetConfigBusinessesSearch()

	q := req.URL.Query()
	q.Add("offset", strconv.Itoa(lgc.Offset))
	q.Add("limit", strconv.Itoa(lgc.Limit))
	q.Add("categories", lgc.Alias)
	q.Add("location", strconv.Itoa(lg.ZipCode))
	q.Add("sort_by", "rating")
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

	fmt.Println("bsr.Error.Code", bsr.Error.Code)
	if bsr.Error.Code != "" {
		log.Fatal(bsr.Error.Code)
	}

	if len(bsr.Businesses) <= 50 {
		lg.IsDoneCategory = true
	}

	var categories = make(map[string]Category)
	var businessCategories []BusinessCategory

	fmt.Println("-----", "Businesses", len(bsr.Businesses))
	for _, b := range bsr.Businesses {
		b.ZipCode = lg.ZipCode

		b.Update()

		for _, c := range b.Categories {
			bc := BusinessCategory{BusinessId: b.Id, CategoryAlias: c.Alias}
			businessCategories = append(businessCategories, bc)
			categories[c.Alias] = c
		}
	}

	fmt.Println("-----", "Categories", len(categories))
	for _, c := range categories {
		c.Update()
	}

	fmt.Println("-----", "BusinessCategories", len(businessCategories))
	for _, bc := range businessCategories {
		bc.Update()
	}

}
