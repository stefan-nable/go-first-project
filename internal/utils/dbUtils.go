package utils

import (
	"Main/internal/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var databaseLink = "localhost:3306"
var databaseName = "log"

func connectToDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp("+databaseLink+")/"+databaseName+"?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func SaveLogsToDB(numWorkers int, resultChannel chan model.Log) {
	db := connectToDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err.Error())
		}
	}(db)

	for i := 0; i < numWorkers; i++ {
		log := <-resultChannel
		model.Log.PrintLog(log)

		_, err := db.Exec("INSERT INTO log (timestamp, message) VALUES (?, ?)",
			log.Timestamp.Format("2006-01-02 15:04:05"),
			log.Message)

		if err != nil {
			panic(err.Error())
		}
	}

}
