package yelp

import (
	"amasia/db"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Log struct {
	Id         int64
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

func (l *Log) Insert() {
	db := db.GetDB()
	now := time.Now()
	res, err := db.Exec(`
		INSERT INTO Log (
			ZipCode,
			Type,
			Config,
			IsDone,
			Error,
			CreatedAt,
			ModifiedAt
		) VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
		;
	`, l.ZipCode, l.Type, l.Config, l.IsDone, l.Error, now, now)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("id", id)

	l.Id = id
	l.CreatedAt = now
	l.ModifiedAt = now
}

func (l Log) Update() {
	db := db.GetDB()
	now := time.Now()
	rows, err := db.Query(`
		UPDATE Log l
		SET
			l.ZipCode=?,
			l.Config=?,
			l.IsDone=?,
			l.Error=?,
			l.ModifiedAt=?
		WHERE l.Id=?
		;
	`, l.ZipCode, l.Config, l.IsDone, l.Error, now, l.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	l.ModifiedAt = now
}

func (l *Log) InitWithLatestBusinessesSearch() {
	db := db.GetDB()
	rows, err := db.Query(`
		SELECT *
    FROM Log l
    ORDER BY IsDone, ModifiedAt DESC
		LIMIT 1
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
	// can i delete the []byte
	err := json.Unmarshal(l.Config, &lc)
	if err != nil {
		log.Fatal(err)
	}

	return lc
}
