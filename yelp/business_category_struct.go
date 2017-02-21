package yelp

import (
	"amasia/db"
	"log"
	"time"
)

type BusinessCategory struct {
	BusinessId    string
	CategoryAlias string
	CreatedAt     time.Time
	ModifiedAt    time.Time
}

func (bc *BusinessCategory) Update() {
	db := db.GetDB()
	now := time.Now()
	rows, err := db.Query(`
		INSERT INTO BusinessCategory (
	    BusinessId,
	    CategoryAlias,
			CreatedAt,
			ModifiedAt
	  ) VALUES (
	    ?,
	    ?,
			?,
	    ?
	  ) ON DUPLICATE KEY UPDATE
	  BusinessId=VALUES(BusinessId),
	  CategoryAlias=VALUES(CategoryAlias),
	  ModifiedAt=VALUES(ModifiedAt)
		;
	`, bc.BusinessId, bc.CategoryAlias, now, now)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	bc.ModifiedAt = now
}
