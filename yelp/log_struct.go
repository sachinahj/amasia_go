package yelp

import (
	"amasia/db"
	"encoding/json"
	"log"
	"time"
)

type Log struct {
	Id         int
	ZipCode    int
	Type       string
	Config     []byte
	IsDone     bool
	Error      string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type LogConfigBusinessesSearch struct {
	Alias  string
	Limit  int
	Offset int
}

func (l Log) Update() {
	db := db.GetDB()
	rows, err := db.Query(`
		update Log l
		set
			l.ZipCode=?,
			l.Config=?,
			l.IsDone=?,
			l.Error=?,
			l.ModifiedAt=?
		where l.Id=?
		;
	`, l.ZipCode, l.Config, l.IsDone, l.Error, time.Now(), l.Id)
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

func (l *Log) InitWithLatestBusinessesSearch() {
	db := db.GetDB()
	rows, err := db.Query(`
		select *
    from Log l
    order by IsDone, CreatedAt desc
		limit 1
		;
	`)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&l.Id, &l.ZipCode, &l.Type, &l.Config, &l.IsDone, &l.Error, &l.CreatedAt, &l.ModifiedAt)
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

func (l Log) GetConfigBusinessesSearch() LogConfigBusinessesSearch {
	var lc LogConfigBusinessesSearch
	err := json.Unmarshal([]byte(l.Config), &lc)
	if err != nil {
		log.Fatal(err)
	}

	return lc
}
