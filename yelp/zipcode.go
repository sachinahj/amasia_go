package yelp

import (
	"amasia/db"
	"log"
)

func GetOldestUpdatedZipCode() ZipCode {
	db := db.GetDB()
	rows, err := db.Query(`
		select zc.*
		from ZipCode zc
		left join
		(
		  select l.Id, l.ZipCode, temp.MaxCreatedAt
		  from Log l
		  inner join
		  (
		    select max(CreatedAt) as MaxCreatedAt, ZipCode
		    from Log l
		    where type="businesses_search"
		    group by l.ZipCode
		  ) temp
		  on l.ZipCode = temp.ZipCode and l.CreatedAt = temp.MaxCreatedAt
		  group by l.ZipCode
		  order by temp.MaxCreatedAt desc
		) temp2
		on zc.ZipCode=temp2.ZipCode
		order by temp2.MaxCreatedAt asc, ZipCode asc
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

// func GetAllZipCodes() []*ZipCode {
// 	db := db.GetDB()
// 	rows, err := db.Query("select * from zipCode")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()
//
// 	zipCodes := make([]*ZipCode, 0)
// 	for rows.Next() {
// 		zc := new(ZipCode)
// 		err := rows.Scan(&zc.ZipCode, &zc.Country, &zc.ForceYelpBusinessSearch, &zc.CreatedAt, &zc.ModifiedAt)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		zipCodes = append(zipCodes, zc)
// 	}
//
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	rows.Close()
// 	return zipCodes
// }
