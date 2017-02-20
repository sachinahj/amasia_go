package yelp

import (
	"amasia/db"
	"amasia/helpers"
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

func (z ZipCode) GetValidCategories() CategoriesConfig {
	var cc CategoriesConfig

	for _, c := range allCategoriesConfig {
		var toKeep = false
		if len(c.CountryWhitelist) == 0 || helpers.StringInSlice(z.Country, c.CountryWhitelist) {
			toKeep = true
		}

		if len(c.CountryBlacklist) == 0 && helpers.StringInSlice(z.Country, c.CountryBlacklist) {
			toKeep = false
		}

		if toKeep {
			cc = append(cc, c)
		}
	}
	return cc
}

func (z *ZipCode) InitWithZipCode() {
	db := db.GetDB()
	rows, err := db.Query(`
		select zc.*
		from ZipCode zc
		where ?=zc.ZipCode
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
