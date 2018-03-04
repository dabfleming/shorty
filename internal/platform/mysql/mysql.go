package mysql

import (
	"database/sql"

	// Load mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Connect connects to the database
func Connect() (*sql.DB, error) {
	// TODO Get this from environment
	url := "username:password@tcp(localhost:3306)/shorty?charset=utf8&parseTime=true"
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
