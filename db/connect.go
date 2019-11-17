package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

// ConnectDB connects to a database
func ConnectDB(host string, port string, user string, password string, dbname string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected!")
	return db, nil
}
