package main

import (
	"fmt"
	"amasia/yelp"
)

type Business struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Rating       float64 `json:"rating"`
	Review_count int64   `json:"review_count"`
	Price        string  `json:"price"`
}

type BusinessesSearchResponse struct {
	Total      int64      `json:"id"`
	Businesses []Business `json:"businesses"`
}

func main() {

	var zc yelp.ZipCode
	lg := yelp.GetLatestLog()

	fmt.Println("lg", lg)
	if lg.ZipCode != 0 {
		zc.ZipCode = lg.ZipCode
		zc.InitWithZipCode()

		fmt.Println("zc", zc)
		filteredCategories := zc.GetValidCategories()
		fmt.Println("filteredCategories", len(filteredCategories))

		if lg.IsDone && lg.Alias ==  filteredCategories[len(filteredCategories) - 1].Alias {
			fmt.Println("set zipcode forced")
			lg.IsDone = false
			lg.Update()
		} else {
			fmt.Println("continue from log")
			lg.IsDone = true
			lg.Update()
		}
	} else {
		fmt.Println("set zipcode forced")
	}



	fmt.Println("\n")

	// zc := yelp.GetOldestUpdatedZipCode()
	// fmt.Println("zc", zc)

	fmt.Println("\n")

	// for _, c := range filteredCategories {
	// 	fmt.Println(c)
	// }


	// viper.SetConfigName("_config") // name of config file (without extension)
	// viper.AddConfigPath(".")       // optionally look for config in the working directory
	// err := viper.ReadInConfig()    // Find and read the config file
	// if err != nil {                // Handle errors reading the config file
	// 	panic(fmt.Errorf("Fatal error config file: %s \n", err))
	// }
	//
	// baseUrl := viper.Get("dataProviders.yelp.baseUrl")
	// rawUrl := fmt.Sprintf("%v/%v", baseUrl, "businesses/search")
	// accessToken := viper.Get("dataProviders.yelp.accessToken")
	//
	// fmt.Println("rawUrl", rawUrl)
	// fmt.Println("accessToken", accessToken)
	//
	// req, err := http.NewRequest("GET", rawUrl, nil)
	// if err != nil {
	// 	panic(fmt.Errorf("Fatal error on request: %s \n", err))
	// }
	//
	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	//
	// q := req.URL.Query()
	// q.Add("location", "78701")
	// q.Add("sort_by", "rating")
	// q.Add("offset", "0")
	// q.Add("limit", "50")
	// q.Add("categories", "yoga")
	// req.URL.RawQuery = q.Encode()
	//
	// fmt.Println("req.URL.String()", req.URL.String())
	// fmt.Println("req.Header", req.Header)
	//
	// client := &http.Client{}
	//
	// resp, err := client.Do(req)
	// fmt.Println("resp", resp)
	//
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	//
	// var s = new(BusinessesSearchResponse)
	// err = json.Unmarshal([]byte(body), &s)
	// if err != nil {
	// 	panic(fmt.Errorf("Fatal error on unmarshall: %s \n", err))
	// }
	//
	// fmt.Println("s", s)

}
