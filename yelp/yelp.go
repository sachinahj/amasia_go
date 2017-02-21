package yelp

import "fmt"

func Start() {
	var zc ZipCode
	zc.InitWithForceBusinessesSearch()
	fmt.Println("zc", zc)

	var lg Log
	if zc.ZipCode != 0 {
		lg = zc.InsertLogBusinessesSearch()
	} else {
		lg.InitWithLatestBusinessesSearch()
	}

	fmt.Println("lg", lg)
	lgc := lg.GetConfigBusinessesSearch()
	fmt.Println("lgc", lgc)

	// Run(lg)
}

func Run(lg Log) {
	fmt.Println("lg", lg)

	var zc ZipCode
	if lg.ZipCode != 0 {
		zc.ZipCode = lg.ZipCode
		zc.InitWithZipCode()

		fmt.Println("zc", zc)
		filteredCategories := zc.GetValidCategories()
		fmt.Println("filteredCategories", len(filteredCategories))

		lgc := lg.GetConfigBusinessesSearch()
		fmt.Println("lgc", lgc)
		fmt.Println("lgc.Alias", lgc.Alias)

		if lg.IsDone && lgc.Alias == filteredCategories[len(filteredCategories)-1].Alias {
			fmt.Println("set zipcode forced")
		} else {
			fmt.Println("continue from log")
			BusinessesSearch(zc, lgc)
		}
	} else {
		fmt.Println("set zipcode forced")
	}
}
