package yelp

import (
	"amasia/db"
	"amasia/helpers"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type ZipCode struct {
	ZipCode               int
	Country               string
	ForceBusinessesSearch bool
	CreatedAt             time.Time
	ModifiedAt            time.Time
}

func (z ZipCode) GetValidCategories() []CategoryConfig {
	var cc []CategoryConfig

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
		SELECT zc.*
		FROM ZipCode zc
		WHERE ?=zc.ZipCode
		LIMIT 1
		;
	`, z.ZipCode)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&z.ZipCode, &z.Country, &z.ForceBusinessesSearch, &z.CreatedAt, &z.ModifiedAt)
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

func (z *ZipCode) InitWithOldestBusinessesSearch() {
	db := db.GetDB()
	rows, err := db.Query(`
		SELECT zc.*
		FROM ZipCode zc
		LEFT JOIN
		(
		  SELECT l.Id, l.ZipCode, temp.MaxModifiedAt
		  FROM Log l
		  INNER JOIN
		  (
		    SELECT max(ModifiedAt) AS MaxModifiedAt, ZipCode
		    FROM Log l
		    WHERE type="businesses_search"
		    GROUP BY l.ZipCode
		  ) temp
		  ON l.ZipCode = temp.ZipCode and l.ModifiedAt = temp.MaxModifiedAt
		  GROUP BY l.ZipCode
		  ORDER BY temp.MaxModifiedAt desc
		) temp2
		ON zc.ZipCode=temp2.ZipCode
		ORDER BY temp2.MaxModifiedAt ASC, ZipCode ASC
		LIMIT 1
		;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&z.ZipCode, &z.Country, &z.ForceBusinessesSearch, &z.CreatedAt, &z.ModifiedAt)
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

func (z *ZipCode) InitWithForceBusinessesSearch() {
	db := db.GetDB()
	rows, err := db.Query(`
		SELECT zc.*
		FROM ZipCode zc
		WHERE ForceBusinessesSearch=TRUE
		LIMIT 1
		;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&z.ZipCode, &z.Country, &z.ForceBusinessesSearch, &z.CreatedAt, &z.ModifiedAt)
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

func (z *ZipCode) InsertLogBusinessesSearch() Log {
	filteredCategories := z.GetValidCategories()
	lgc := LogConfigBusinessesSearch{Alias: filteredCategories[0].Alias, Limit: 50, Offset: 0}
	fmt.Println("lgc", lgc)
	lgc_byte, err := json.Marshal(lgc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("lgc_byte", lgc_byte)
	lg := Log{ZipCode: z.ZipCode, Type: "businesses_search", Config: lgc_byte, IsDone: false, Error: ""}
	lg.Insert()
	lg.Update()
	return lg
}
