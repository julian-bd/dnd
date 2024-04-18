package data

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() error {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "dnd",
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	return db.Ping()
}

func HasBeenSeeded() bool {
	row := db.QueryRow("SELECT id FROM playable_race limit 1")
	var id int
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	return true
}
