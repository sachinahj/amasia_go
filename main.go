package main

import (
	"amasia/config"
	"amasia/yelp"
	"fmt"
)

func main() {

	var environment = config.Get("environment")
	fmt.Println("environment", environment)

	if environment == "development" {
		var zc yelp.ZipCode
		zc.ZipCode = 10027

		fmt.Println("zc.ZipCode", zc.ZipCode)
		zc.RunAnalysis()

	} else {
		yelp.Start()
	}

}
