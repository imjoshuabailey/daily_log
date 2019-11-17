package main

import (
	"log"
	"os"

	"github.com/imjoshuabailey/arnold_automator/db"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	var (
		host     = os.Getenv(`POSTGRES_HOST`)
		port     = os.Getenv(`POSTGRES_PORT`)
		user     = os.Getenv(`POSTGRES_USER`)
		password = os.Getenv(`POSTGRES_PASSWORD`)
		dbname   = os.Getenv(`POSTGRES_DB`)
	)

	conn, err := db.ConnectDB(host, port, user, password, dbname)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer conn.Close()

	err = db.CreateTable(conn)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

}
