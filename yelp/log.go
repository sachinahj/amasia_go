package yelp

import (
	"amasia/db"
	"log"
)

func GetLatestLog() Log {
	db := db.GetDB()
	rows, err := db.Query(`
		select *
    from yelpLogBusinessSearch l
    order by isDone, modifiedAt desc
		limit 1
		;
	`)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var lg Log

	for rows.Next() {
		err := rows.Scan(&lg.Id, &lg.ZipCode, &lg.Alias, &lg.Limit, &lg.Offset, &lg.IsDone, &lg.Error, &lg.CreatedAt, &lg.ModifiedAt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	return lg
}