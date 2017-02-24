package yelp

import (
	"amasia/db"
	"amasia/helpers"
	"log"
	"time"
)

type ZipCode struct {
	ZipCode               int
	City                  string
	State                 string
	Country               string
	ForceBusinessesSearch bool
	CreatedAt             time.Time
	ModifiedAt            time.Time
}

func (z ZipCode) RunAnalysis() {
	db := db.GetDB()

	_, err := db.Exec(`
		INSERT INTO amasia.ZipCodeCategoriesLevel4
		(
      ZipCode,
      Alias,
      Count
    )
		(
			SELECT
			  b.LocationZipCode as ZipCode,
			  yct.AliasLevel4 as Alias,
			  count(yct.AliasLevel4) as Count

			FROM Business b

			JOIN BusinessCategory bc ON b.Id=bc.BusinessId
			JOIN Category yc ON yc.Alias=bc.CategoryAlias
			JOIN CategoryTree yct ON yct.AliasLevel4=yc.Alias

			AND b.LocationZipCode=?

			GROUP BY b.LocationZipCode, yct.AliasLevel4
			ORDER BY b.LocationZipCode, Count DESC
		) ON DUPLICATE KEY UPDATE
		Count=Values(Count)
		;
		INSERT INTO amasia.ZipCodeCategoriesLevel3
		(
      ZipCode,
      Alias,
      Count
    )
		(
			SELECT
			  b.LocationZipCode as ZipCode,
			  yct.AliasLevel3 as Alias,
			  count(yct.AliasLevel3) as Count

			FROM Business b

			JOIN BusinessCategory bc ON b.Id=bc.BusinessId
			JOIN Category yc ON yc.Alias=bc.CategoryAlias
			JOIN CategoryTree yct ON yct.AliasLevel3=yc.Alias

			AND b.LocationZipCode=?

			GROUP BY b.LocationZipCode, yct.AliasLevel3
			ORDER BY b.LocationZipCode, Count DESC
		) ON DUPLICATE KEY UPDATE
		Count=Values(Count)
		;
		INSERT INTO amasia.ZipCodeCategoriesLevel2
		(
      ZipCode,
      Alias,
      Count
    )
		(
			SELECT
			  b.LocationZipCode as ZipCode,
			  yct.AliasLevel2 as Alias,
			  count(yct.AliasLevel2) as Count

			FROM Business b

			JOIN BusinessCategory bc ON b.Id=bc.BusinessId
			JOIN Category yc ON yc.Alias=bc.CategoryAlias
			JOIN CategoryTree yct ON yct.AliasLevel2=yc.Alias

			AND b.LocationZipCode=?

			GROUP BY b.LocationZipCode, yct.AliasLevel2
			ORDER BY b.LocationZipCode, Count DESC
		) ON DUPLICATE KEY UPDATE
		Count=Values(Count)
		;
		INSERT INTO amasia.ZipCodeCategoriesLevel1
		(
      ZipCode,
      Alias,
      Count
    )
		(
			SELECT
			  b.LocationZipCode as ZipCode,
			  yct.AliasLevel1 as Alias,
			  count(yct.AliasLevel1) as Count

			FROM Business b

			JOIN BusinessCategory bc ON b.Id=bc.BusinessId
			JOIN Category yc ON yc.Alias=bc.CategoryAlias
			JOIN CategoryTree yct ON yct.AliasLevel1=yc.Alias

			AND b.LocationZipCode=?

			GROUP BY b.LocationZipCode, yct.AliasLevel1
			ORDER BY b.LocationZipCode, Count DESC
		) ON DUPLICATE KEY UPDATE
		Count=Values(Count)
		;
	`, z.ZipCode, z.ZipCode, z.ZipCode, z.ZipCode)

	if err != nil {
		log.Fatal(err)
	}

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
		err := rows.Scan(&z.ZipCode, &z.City, &z.State, &z.Country, &z.ForceBusinessesSearch, &z.CreatedAt, &z.ModifiedAt)
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
		err := rows.Scan(&z.ZipCode, &z.City, &z.State, &z.Country, &z.ForceBusinessesSearch, &z.CreatedAt, &z.ModifiedAt)
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
		err := rows.Scan(&z.ZipCode, &z.City, &z.State, &z.Country, &z.ForceBusinessesSearch, &z.CreatedAt, &z.ModifiedAt)
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
