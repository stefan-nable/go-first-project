package memory

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var databaseLink = "localhost:3306"
var databaseName = "log"

func ConnectToDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp("+databaseLink+")/"+databaseName+"?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func FlushDB(db *sql.DB) {
	_, err := db.Exec("DELETE FROM log")
	if err != nil {
		panic(err.Error())
	}
}
