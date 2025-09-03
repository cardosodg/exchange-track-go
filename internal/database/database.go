package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type HolidayDB struct {
	ID   int    `db:"id"`
	Date string `db:"date"`
	Name string `db:"name"`
	Type string `db:"type"`
	Year string `db:"year"`
}

type ExchangeDataDB struct {
	ID         int     `db:"id"`
	Bid        float64 `db:"bid"`
	Timestamp  string  `db:"timestamp"`
	CreateDate string  `db:"create_date"`
	Type       string  `db:"type"`
}

type HistoricalExchange struct {
	ID       int     `db:"id"`
	Day      string  `db:"day"`
	AvgBid   float64 `db:"avg_bid"`
	FinalBid float64 `db:"final_bid"`
	Type     string  `db:"type"`
}

func Connect(user, pass, host, port, dbName string) *sql.DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbName)

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

func CreateTables(sql *sql.DB) bool {
	fmt.Println("Table created")
	return true
}

func InsertExchangeData(db *sql.DB, rate ExchangeDataDB) error {
	query := `INSERT INTO exchange_rates (bid, timestamp, create_date, type) 
			  VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, rate.Bid, rate.Timestamp, rate.CreateDate, rate.Type)
	if err != nil {
		return fmt.Errorf("erro ao inserir dados: %v", err)
	}
	return nil
}

func InserAlltHolidayData(db *sql.DB, holidays []HolidayDB) error {
	query := `INSERT INTO holidays (date, name, type, year) VALUES `

	var values []any

	for i, holiday := range holidays {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?, ?)"

		values = append(values, holiday.Date, holiday.Name, holiday.Type, holiday.Year)
	}

	_, err := db.Exec(query, values)
	if err != nil {
		return fmt.Errorf("error inserting data: %v", err)
	}
	return nil
}
