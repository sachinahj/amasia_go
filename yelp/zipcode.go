package yelp

import (
	"amasia/db"
	"fmt"
	"log"
	"time"
)

type ZipCode struct {
	ZipCode                 int
	Country                 string
	ForceYelpBusinessSearch bool
	CreatedAt               time.Time
	ModifiedAt              time.Time
}

func init() {
	fmt.Println("initing yelp zipcode")
}

func GetAllZipCodes() []*ZipCode{
	fmt.Println("Getting all zipcodes")

	db := db.GetDB()
	rows, err := db.Query("select * from zipCode")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	zipCodes := make([]*ZipCode, 0)
	for rows.Next() {
		zc := new(ZipCode)
		err := rows.Scan(&zc.ZipCode, &zc.Country, &zc.ForceYelpBusinessSearch, &zc.CreatedAt, &zc.ModifiedAt)
		if err != nil {
			log.Fatal(err)
		}
		zipCodes = append(zipCodes, zc)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	return zipCodes;
}
