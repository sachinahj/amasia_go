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

func (z *ZipCode) getValidCategories() {
	fmt.Println("getValidCategories for", z)
}

func (z *ZipCode) GetWithZipCode() {
	db := db.GetDB()
	rows, err := db.Query(`
		select zc.*
		from zipCode zc
		where ?=zc.zipCode
		limit 1
		;
	`, z.ZipCode)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&z.ZipCode, &z.Country, &z.ForceYelpBusinessSearch, &z.CreatedAt, &z.ModifiedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
}

func init() {}

func GetOldestUpdatedZipCode() ZipCode {
	db := db.GetDB()
	rows, err := db.Query(`
		select zc.*
		from zipCode zc
		left join
		(
			select l.id, l.zipCode, l.alias, temp.maxModifiedAt
			from yelpLogBusinessSearch l
			inner join
			(
				select max(modifiedAt) as maxModifiedAt, zipCode
				from yelpLogBusinessSearch l
				group by l.zipCode
			) temp
			on l.zipCode = temp.zipCode and l.modifiedAt = temp.maxModifiedAt
			group by l.zipCode
			order by temp.maxModifiedAt desc
		) temp2
		on zc.zipCode=temp2.zipCode
		order by temp2.maxModifiedAt asc, zipCode asc
		limit 1
		;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var zc ZipCode

	for rows.Next() {
		err := rows.Scan(&zc.ZipCode, &zc.Country, &zc.ForceYelpBusinessSearch, &zc.CreatedAt, &zc.ModifiedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	return zc
}

func GetAllZipCodes() []*ZipCode {
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
	return zipCodes
}
