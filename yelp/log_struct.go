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

type LogConfigZipCodeAnalysis struct {
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

func (l *Log) Update() {
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

func (l *Log) InitWithLatestBusinessesSearch() {
	db := db.GetDB()
	rows, err := db.Query(`
		SELECT *
    FROM Log l
		WHERE Type="businesses_search"
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

func (l *Log) InitWithNotDoneZipCodeAnalysis() {
	db := db.GetDB()
	rows, err := db.Query(`
		SELECT *
    FROM Log l
		WHERE Type="zip_code_analysis"
		AND IsDone=false
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
	nextType := "businesses_search"
	var zc = ZipCode{ZipCode: l.ZipCode}
	zc.InitWithZipCode()

	var lgcbs LogConfigBusinessesSearch
	var lgcza LogConfigZipCodeAnalysis

	switch l.Type {
	case "businesses_search":
		lgcbs = l.GetConfigBusinessesSearch()
		filteredCategories := zc.GetValidCategories()

		var i int
		var c CategoryConfig
		if l.IsDoneCategory {

			for i, c = range filteredCategories {
				if c.Alias == lgcbs.Alias {
					break
				}
			}

			if i+1 < len(filteredCategories) {
				lgcbs = LogConfigBusinessesSearch{Alias: filteredCategories[i+1].Alias, Limit: 50, Offset: 0}
			} else {
				nextType = "zip_code_analysis"
				lgcza = LogConfigZipCodeAnalysis{}
			}

		} else {
			lgcbs = LogConfigBusinessesSearch{Alias: lgcbs.Alias, Limit: 50, Offset: (lgcbs.Offset + 50)}
		}
	case "zip_code_analysis":
	default:
		nextType = "businesses_search"
		zc.InitWithForceBusinessesSearch()
		if zc.ZipCode == 0 {
			zc.InitWithOldestBusinessesSearch()
		} else {
			zc.ForceBusinessesSearch = false
			zc.Update()
		}
		l.ZipCode = zc.ZipCode
		l.InitWithNewBusinessesSearch()
		lgcbs = l.GetConfigBusinessesSearch()
	}

	l.Id = 0
	l.ZipCode = zc.ZipCode
	l.Type = nextType

	switch nextType {
	case "businesses_search":
		lgcbs_byte, err := json.Marshal(lgcbs)
		if err != nil {
			log.Fatal(err)
		}
		l.Config = lgcbs_byte
	case "zip_code_analysis":
		lgcza_byte, err := json.Marshal(lgcza)
		if err != nil {
			log.Fatal(err)
		}
		l.Config = lgcza_byte
	}

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

func (l Log) GetConfigZipCodeAnalysis() LogConfigZipCodeAnalysis {
	var lc LogConfigZipCodeAnalysis
	// can i delete the []byte
	err := json.Unmarshal(l.Config, &lc)
	if err != nil {
		log.Fatal(err)
	}

	return lc
}
