package db

import (
	"fmt"
)

func GetAllZipCodes() {
	fmt.Println("db: Getting all zipcodes")

	runQuery("select * from zipCode")

	// fmt.Println(results)
}
