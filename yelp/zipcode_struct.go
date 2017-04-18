package yelp

import (
	"amasia/config"
	"amasia/db"
	"amasia/helpers"
	"fmt"
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

func (z *ZipCode) Update() {
	db := db.GetDB()
	now := time.Now()
	rows, err := db.Query(`
		UPDATE ZipCode z
		SET
			z.City=?,
			z.State=?,
			z.Country=?,
			z.ForceBusinessesSearch=?,
			z.ModifiedAt=?
		WHERE z.ZipCode=?
		;
	`, z.City, z.State, z.Country, z.ForceBusinessesSearch, now, z.ZipCode)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	z.ModifiedAt = now
}

func (z ZipCode) RunAnalysis() {
	db := db.GetDB()
	now := time.Now()
	var err error
	var query string

	var dbName = config.Get("database.main.database")

	// transfer undocumented zip codes
	_, err = db.Exec(`
		INSERT INTO ZipCode
		(
		  ZipCode,
		  City,
		  State,
		  Country,
		  ForceBusinessesSearch,
			CreatedAt,
			ModifiedAt
		)
		(
		  SELECT
		    b.LocationZipCode as ZipCode,
		    b.LocationCity as City,
		    b.LocationState as State,
		    b.LocationCountry as Country,
		    0 as ForceBusinessesSearch,
		    ? as CreatedAt,
		    ? as ModifiedAt
		  FROM Business b
		  LEFT JOIN ZipCode z on z.ZipCode = b.LocationZipCode
		  WHERE z.ZipCode IS NULL
		  AND b.LocationZipCode > 10000
		  AND b.LocationZipCode < 99999
		  GROUP BY b.LocationZipCode
		)
		;
	`, now, now)

	if err != nil {
		log.Fatal(err)
	}

	// insert updated analysis into each ZipCodeCategoriesLevel table
	for i := 1; i <= 4; i++ {

		query = fmt.Sprintf(`
			INSERT INTO %s.ZipCodeCategoriesLevel%d
			(
		    ZipCode,
		    Alias,
		    Count,
				AverageRating,
				AverageReviewCount,
				AveragePrice
		  )
			(
				SELECT
				  b.LocationZipCode as ZipCode,
				  yct.AliasLevel%d as Alias,
				  count(yct.AliasLevel%d) as Count,
					avg(b.Rating) as AverageRating,
			    avg(b.ReviewCount) as AverageReviewCount,
			    avg(b.Price) as AveragePrice

				FROM Business b

				JOIN BusinessCategory bc ON b.Id=bc.BusinessId
				JOIN Category yc ON yc.Alias=bc.CategoryAlias
				JOIN CategoryTree yct ON yct.AliasLevel4=yc.Alias

				AND b.LocationZipCode=%d

				GROUP BY b.LocationZipCode, yct.AliasLevel%d
			) ON DUPLICATE KEY UPDATE
			Count=Values(Count)
		`, dbName, i, i, i, z.ZipCode, i)

		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	// delete all BusinessCategory
	_, err = db.Exec(`
		DELETE BusinessCategory
		FROM BusinessCategory
		JOIN Business ON Business.Id = BusinessCategory.BusinessId
		WHERE Business.ZipCode=?
		;
	`, z.ZipCode)

	if err != nil {
		log.Fatal(err)
	}

	// delete all Business
	_, err = db.Exec(`
		DELETE Business
		FROM Business
		WHERE Business.ZipCode=?
	`, z.ZipCode)

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
