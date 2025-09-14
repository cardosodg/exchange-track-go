package database

import (
	"ExchangeTrack/internal/config"
	"ExchangeTrack/internal/model"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	config := config.LoadConfig()

	dnsRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort)

	rootDB, err := sql.Open("mysql", dnsRoot)

	if err != nil {
		log.Fatal("Error connecting to MariaDB", err)
	}
	defer rootDB.Close()

	_, err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.DBName))
	if err != nil {
		log.Fatal("Error creating database:", err)
	}
	log.Printf("Database %s verified/created\n", config.DBName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database inaccessible:", err)
	}

	log.Println("Successfully connected to the database!")
	return db
}

func Close(db *sql.DB) {
	if db != nil {
		db.Close()
		log.Println("Database connection closed.")
	}
}

func IsTableEmpty(sql *sql.DB, tableName string) bool {
	fmt.Println(tableName)
	return true
}

func CreateTables(db *sql.DB) {
	queryRates := `
	CREATE TABLE IF NOT EXISTS exchange_rates (
		id INT AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(10) NOT NULL,
		timestamp VARCHAR(50),
		create_date VARCHAR(50),
		bid DOUBLE,
		high DOUBLE,
		low DOUBLE,
		average DOUBLE
	);`

	queryHist := `
	CREATE TABLE IF NOT EXISTS exchange_hist (
		id INT AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(10) NOT NULL,
		timestamp VARCHAR(50),
		create_date VARCHAR(50),
		bid DOUBLE,
		high DOUBLE,
		low DOUBLE,
		average DOUBLE
	);`

	_, err := db.Exec(queryRates)
	if err != nil {
		log.Fatal("Error in creating table exchange_rates", err)
	}

	_, err = db.Exec(queryHist)
	if err != nil {
		log.Fatal("Error in creating table exchange_hist", err)
	}

	log.Println("Tables exchange_rates and exchange_hist created.")
}

func InsertExchangeData(db *sql.DB, table string, rates []model.CurrencyData) error {
	query := fmt.Sprintf(`INSERT INTO %s (code, timestamp, create_date, bid, high, low, average) VALUES (?, ?, ?, ?, ?, ?, ?)`, table)

	for _, rate := range rates {
		_, err := db.Exec(query, rate.Code, rate.Timestamp, rate.CreateDate, rate.Bid, rate.High, rate.Low, rate.Average)
		if err != nil {
			return fmt.Errorf("error inserting %s: %v", rate.Code, err)
		}
	}

	return nil
}

// func InserAlltHolidayData(db *sql.DB, holidays []HolidayDB) error {
// 	query := `INSERT INTO holidays (date, name, type, year) VALUES `

// 	var values []any

// 	for i, holiday := range holidays {
// 		if i > 0 {
// 			query += ", "
// 		}
// 		query += "(?, ?, ?, ?)"

// 		values = append(values, holiday.Date, holiday.Name, holiday.Type, holiday.Year)
// 	}

// 	_, err := db.Exec(query, values)
// 	if err != nil {
// 		return fmt.Errorf("error inserting data: %v", err)
// 	}
// 	return nil
// }
