package yelp

import (
	"amasia/db"
	"log"
)

func GetOldestUpdatedZipCode() ZipCode {
	db := db.GetDB()
	rows, err := db.Query(`
		SELECT zc.*
		FROM ZipCode zc
		LEFT JOIN
		(
		  SELECT l.Id, l.ZipCode, temp.MaxCreatedAt
		  FROM Log l
		  INNER JOIN
		  (
		    SELECT max(CreatedAt) AS MaxCreatedAt, ZipCode
		    FROM Log l
		    WHERE type="businesses_search"
		    GROUP BY l.ZipCode
		  ) temp
		  ON l.ZipCode = temp.ZipCode and l.CreatedAt = temp.MaxCreatedAt
		  GROUP BY l.ZipCode
		  ORDER BY temp.MaxCreatedAt desc
		) temp2
		ON zc.ZipCode=temp2.ZipCode
		ORDER BY temp2.MaxCreatedAt ASC, ZipCode ASC
		LIMIT 1
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
