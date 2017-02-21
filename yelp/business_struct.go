package yelp

import (
	"amasia/db"
	"log"
	"time"
)

type Business struct {
	Id          string `json:"id"`
	ZipCode     int
	Name        string     `json:"name"`
	Rating      float64    `json:"rating"`
	ReviewCount int64      `json:"review_count"`
	Price       string     `json:"price"`
	Categories  []Category `json:"categories"`
	Location    struct {
		ZipCode string `json:"zip_code"`
		City    string `json:"city"`
		State   string `json:"state"`
		Country string `json:"country"`
	} `json:"location"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func (b *Business) Update() {
	db := db.GetDB()
	now := time.Now()
	price := len(b.Price)
	rows, err := db.Query(`
		INSERT INTO Business (
      Id,
      ZipCode,
      Name,
      Rating,
      ReviewCount,
      Price,
      LocationZipCode,
      LocationCity,
      LocationState,
      LocationCountry,
      CoordinatesLatitude,
      CoordinatesLongitude,
      CreatedAt,
      ModifiedAt
    ) VALUES (
      ?,
      ?,
      ?,
      ?,
      ?,
      ?,
      ?,
      ?,
      ?,
      ?,
      ?,
      ?,
      ?,
      ?
    ) ON DUPLICATE KEY UPDATE
    Id=VALUES(Id),
    ZipCode=VALUES(ZipCode),
    Name=VALUES(Name),
    Rating=VALUES(Rating),
    ReviewCount=VALUES(ReviewCount),
    Price=VALUES(Price),
    LocationZipCode=VALUES(LocationZipCode),
    LocationCity=VALUES(LocationCity),
    LocationState=VALUES(LocationState),
    LocationCountry=VALUES(LocationCountry),
    CoordinatesLatitude=VALUES(CoordinatesLatitude),
    CoordinatesLongitude=VALUES(CoordinatesLongitude),
    ModifiedAt=VALUES(ModifiedAt)
		;
	`, b.Id, b.ZipCode, b.Name, b.Rating, b.ReviewCount, price, b.Location.ZipCode, b.Location.City, b.Location.State, b.Location.Country, b.Coordinates.Latitude, b.Coordinates.Longitude, now, now)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows.Close()
	b.ModifiedAt = now
}
