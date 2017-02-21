package yelp

import (
	"amasia/db"
	"encoding/json"
	"log"
	"time"
)

type Log struct {
	Id             int64
	ZipCode        int
	Type           string
	Config         []byte
	IsDone         bool
	Error          string
	IsDoneCategory bool
	CreatedAt      time.Time
	ModifiedAt     time.Time
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
    ORDER BY IsDone ASC, ModifiedAt DESC
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

func (l *Log) InitWithNextLog() {
	var zc = ZipCode{ZipCode: l.ZipCode}

	zc.InitWithZipCode()
	lgc := l.GetConfigBusinessesSearch()
	filteredCategories := zc.GetValidCategories()

	var i int
	var c CategoryConfig
	if l.IsDoneCategory {
		for i, c = range filteredCategories {
			if c.Alias == lgc.Alias {
				break
			}
		}

		lgc = LogConfigBusinessesSearch{Alias: filteredCategories[i+1].Alias, Limit: 50, Offset: 0}
	} else {
		lgc = LogConfigBusinessesSearch{Alias: lgc.Alias, Limit: 50, Offset: (lgc.Offset + 50)}
	}

	lgc_byte, err := json.Marshal(lgc)
	if err != nil {
		log.Fatal(err)
	}

	l.Id = 0
	l.ZipCode = zc.ZipCode
	l.Type = "businesses_search"
	l.Config = lgc_byte
	l.IsDone = false
	l.Error = ""
}

func (l *Log) InitWithNewBusinessesSearch() {
	var zc = ZipCode{ZipCode: l.ZipCode}
	zc.InitWithZipCode()
	filteredCategories := zc.GetValidCategories()

	lgc := LogConfigBusinessesSearch{Alias: filteredCategories[0].Alias, Limit: 50, Offset: 0}
	lgc_byte, err := json.Marshal(lgc)
	if err != nil {
		log.Fatal(err)
	}

	l.Id = 0
	l.ZipCode = zc.ZipCode
	l.Type = "businesses_search"
	l.Config = lgc_byte
	l.IsDone = false
	l.Error = ""
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
