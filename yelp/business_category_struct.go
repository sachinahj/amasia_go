package yelp

import "fmt"

type BusinessCategory struct {
	BusinessId    string
	CategoryAlias string
}

func (bc *BusinessCategory) Update() {

	fmt.Println("update this shit", bc)
	// db := db.GetDB()
	//
	// now := time.Now()
	// price := len(b.Price)
	//
	// rows, err := db.Query(`
	// 	INSERT INTO Business
	// 	(
	//     Id,
	//     ZipCode,
	//     Name,
	//     Rating,
	//     ReviewCount,
	//     Price,
	//     LocationZipCode,
	//     LocationCity,
	//     LocationState,
	//     LocationCountry,
	//     CoordinatesLatitude,
	//     CoordinatesLongitude,
	//     CreatedAt,
	//     ModifiedAt
	//   ) VALUES
	//   (
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?,
	//     ?
	//   ) ON DUPLICATE KEY UPDATE
	//   Id=VALUES(Id),
	//   ZipCode=VALUES(ZipCode),
	//   Name=VALUES(Name),
	//   Rating=VALUES(Rating),
	//   ReviewCount=VALUES(ReviewCount),
	//   Price=VALUES(Price),
	//   LocationZipCode=VALUES(LocationZipCode),
	//   LocationCity=VALUES(LocationCity),
	//   LocationState=VALUES(LocationState),
	//   LocationCountry=VALUES(LocationCountry),
	//   CoordinatesLatitude=VALUES(CoordinatesLatitude),
	//   CoordinatesLongitude=VALUES(CoordinatesLongitude),
	//   ModifiedAt=VALUES(ModifiedAt)
	// 	;
	// `, b.Id, b.ZipCode, b.Name, b.Rating, b.ReviewCount, price, b.Location.ZipCode, b.Location.City, b.Location.State, b.Location.Country, b.Coordinates.Latitude, b.Coordinates.Longitude, now, now)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	//
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// rows.Close()
	// b.ModifiedAt = now
}
