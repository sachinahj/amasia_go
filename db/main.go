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
