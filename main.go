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
		yelp.Sandbox()
	} else {
		yelp.Start()
	}

}
