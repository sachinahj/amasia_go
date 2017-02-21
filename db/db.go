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
	host := config.Get("database.yelp.host")
	port := config.Get("database.yelp.port")
	username := config.Get("database.yelp.username")
	password := config.Get("database.yelp.password")
	database := config.Get("database.yelp.database")

	connectionUrl := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", username, password, host, port, database)
	fmt.Println("db: connectionUrl", connectionUrl)

	var err error
	db, err = sql.Open("mysql", connectionUrl)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {

}

func GetDB() *sql.DB {
	return db;
}
