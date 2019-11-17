package db

import (
	"database/sql"
	"log"
)

// CreateTable creates a table
func CreateTable(db *sql.DB) error {
	const qry = `
	CREATE TABLE IF NOT EXISTS vehicles (
		vin VARCHAR(17),
		unit_number VARCHAR(8),
		class VARCHAR(8),
		year INT,
		model VARCHAR(8),
		body_style VARCHAR(8),
		color VARCHAR(16),
		license_number VARCHAR(8),
		miles INT,
		status VARCHAR(18),
		contract_number INT,
		date_last_used VARCHAR(10),
		stock_number VARCHAR(8)
	)`

	result, err := db.Exec(qry)
	if err != nil {
		return err
	}

	log.Printf("Table created! - %+v", result)

	return nil
}
