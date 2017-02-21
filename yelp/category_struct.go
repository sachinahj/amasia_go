package yelp

import (
	"amasia/db"
	"log"
	"time"
)

type CategoryConfig struct {
	Alias            string   `json:"alias"`
	Title            string   `json:"title"`
	CountryWhitelist []string `json:"country_whitelist"`
	CountryBlacklist []string `json:"country_blacklist"`
	Parents          []string `json:"parents"`
}

type Category struct {
	Alias      string `json:"alias"`
	Title      string `json:"title"`
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func (c *Category) Update() {

	db := db.GetDB()

	now := time.Now()

	rows, err := db.Query(`
		INSERT INTO Category
		(
	    Alias,
	    Title,
			CreatedAt,
			ModifiedAt
	  ) VALUES
	  (
	    ?,
	    ?,
			?,
	    ?
	  ) ON DUPLICATE KEY UPDATE
	  Alias=VALUES(Alias),
	  Title=VALUES(Title),
	  ModifiedAt=VALUES(ModifiedAt)
		;
	`, c.Alias, c.Title, now, now)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	c.ModifiedAt = now
}
