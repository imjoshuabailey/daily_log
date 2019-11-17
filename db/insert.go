package db

import (
	"database/sql"
	"log"

	"github.com/imjoshuabailey/arnold_automator/lib"
)

// Insert inserts information into a database
func Insert(db *sql.DB, v lib.Vehicle) error {
	sqlStatement := `
	INSERT INTO vehicles (
		vin, 
		unit_number, 
		class, 
		year, 
		model, 
		body_style, 
		color, 
		license_number, 
		miles, 
		status, 
		contract_number,
		date_last_used, 
		stock_number
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	result, err := db.Exec(
		sqlStatement,
		v.VIN,
		v.UnitNumber,
		v.Class,
		v.Year,
		v.Model,
		v.BodyStyle,
		v.Color,
		v.LicenseNumber,
		v.Miles,
		v.Status,
		v.ContractNumber,
		v.DateLastUsed,
		v.StockNumber,
	)
	if err != nil {
		return err
	}

	log.Printf("sucessfully inserted! - %+v", result)
	return nil
}
