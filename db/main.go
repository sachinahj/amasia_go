package db

import (
	"amasia/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	host := config.Get("database.main.host")
	port := config.Get("database.main.port")
	username := config.Get("database.main.username")
	password := config.Get("database.main.password")
	database := config.Get("database.main.database")

	connectionUrl := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", username, password, host, port, database)
	fmt.Println("db: connectionUrl", connectionUrl)

	var err error
	db, err = sql.Open("mysql", connectionUrl)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

}

func runQuery(query string) {
	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

  columns, _ := rows.Columns()
  count := len(columns)
  values := make([]interface{}, count)
  valuePtrs := make([]interface{}, count)

	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}

  for rows.Next() {
		fmt.Println(columns)
		fmt.Println(valuePtrs)

    rows.Scan(valuePtrs...)

    for i, col := range columns {

      var v interface{}

      val := values[i]

      b, ok := val.([]byte)

      if ok {
          v = string(b)
      } else {
          v = val
      }

			fmt.Println(col, v)

    }
  }

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	rows.Close()

	// return results
}
