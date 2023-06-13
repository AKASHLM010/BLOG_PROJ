package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectToDB() error {
	// Connection parameters
	dbHost := "localhost"
	dbPort := 5432
	dbUser     := "postgres"
	dbPassword := "723101"
	dbName     := "blogb"
	

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		return err
	}

	DB = db
	return nil
}
