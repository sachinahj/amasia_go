package main

import (
	"amasia/yelp"
)

func main() {

	yelp.Start()

	// var zc yelp.ZipCode
	// var lg yelp.Log
	// lg.InitWithLatestBusinessesSearch()
	//
	// fmt.Println("lg", lg)
	// if lg.ZipCode != 0 {
	// 	zc.ZipCode = lg.ZipCode
	// 	zc.InitWithZipCode()
	//
	// 	fmt.Println("zc", zc)
	// 	filteredCategories := zc.GetValidCategories()
	// 	fmt.Println("filteredCategories", len(filteredCategories))
	//
	// 	lgc := lg.GetConfigBusinessesSearch()
	// 	fmt.Println("lgc", lgc)
	// 	fmt.Println("lgc.Alias", lgc.Alias)
	//
	// 	if lg.IsDone && lgc.Alias ==  filteredCategories[len(filteredCategories) - 1].Alias {
	// 		fmt.Println("set zipcode forced")
	// 	} else {
	// 		fmt.Println("continue from log")
	// 		yelp.BusinessesSearch(zc, lgc)
	// 	}
	// } else {
	// 	fmt.Println("set zipcode forced")
	// }
	//
	//
	//
	// fmt.Println("\n")
	//
	// // zc := yelp.GetOldestUpdatedZipCode()
	// // fmt.Println("zc", zc)
	//
	// fmt.Println("\n")
	//
	// // for _, c := range filteredCategories {
	// // 	fmt.Println(c)
	// // }

}
