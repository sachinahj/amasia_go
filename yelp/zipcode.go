package yelp

import (
	"amasia/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type ZipCode struct {
	ZipCode                 string
	Country                 []string
	ForceYelpBusinessSearch bool
}

func init() {
	fmt.Println("initing yelp zipcode")

	host := config.Get("database.main.host")
	port := config.Get("database.main.port")
	username := config.Get("database.main.username")
	password := config.Get("database.main.password")
	database := config.Get("database.main.database")

	connectionUrl := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)
	fmt.Println("connectionUrl", connectionUrl)

	db, err := sql.Open("mysql", connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from zipCode")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()


  columns, _ := rows.Columns()
  count := len(columns)
  values := make([]interface{}, count)
  valuePtrs := make([]interface{}, count)

	var allZipCodes []ZipCode;

  for rows.Next() {

      for i, _ := range columns {
          valuePtrs[i] = &values[i]
      }

      rows.Scan(valuePtrs...)

      for i, col := range columns {

          var v interface{}

          val := values[i]

          b, ok := val.([]byte)

          if (ok) {
              v = string(b)
          } else {
              v = val
          }

					fmt.Println(col, v)
					zipCode := v.(ZipCode)
					allZipCodes = append(allZipCodes, zipCode)

      }
  }

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(allZipCodes)

}
