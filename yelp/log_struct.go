package yelp

import (
	"amasia/db"
	"log"
	"time"
)

type Log struct {
	Id         int
	ZipCode    int
	Alias      string
	Limit      int
	Offset     int
	IsDone     bool
	Error      string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func (l Log) Update() {
	db := db.GetDB()
	rows, err := db.Query(`
		update YelpLogBusinessSearch l
		set
			l.ZipCode=?,
			l.Alias=?,
			l.Limit=?,
			l.Offset=?,
			l.IsDone=?,
			l.Error=?,
			l.ModifiedAt=?
		where l.Id=?
		;
	`, l.ZipCode, l.Alias, l.Limit, l.Offset, l.IsDone, l.Error, l.ModifiedAt, l.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
}
