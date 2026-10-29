package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_apigolang")
	if err != nil {
		panic("Gagal Tersambung ke server!")
	}
	if err = db.Ping(); err != nil {
		panic("Gagal Tersambung ke Database: " + err.Error())
	}
	return db
}
