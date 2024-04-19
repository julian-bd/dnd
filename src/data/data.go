package data

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/gofor-little/env"
)

var db *sql.DB

func InitDB() error {
    user, err := env.MustGet("DB_USER")
    if err != nil {
        return err
    }
    pass, err := env.MustGet("DB_PASS")
    if err != nil {
        return err
    }
	cfg := mysql.Config{
		User:   user,
		Passwd: pass,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "dnd",
        MultiStatements: true,
	}
    db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
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

func CreateTables() error {
    // TODO: is this the right path?
    filePath := filepath.Join("data/seed.sql")
    file, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    sql := string(file)
    _, err = db.Exec(sql)
    if err != nil {
        return err
    }
    return nil
}
