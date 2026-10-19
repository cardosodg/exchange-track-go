package main

import (
	"ExchangeTrack/internal/config"
	"ExchangeTrack/internal/database"
	"ExchangeTrack/internal/datetime"
	"ExchangeTrack/internal/service"
	"database/sql"
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

func checkNonWorkingDay() {
	now := time.Now()

	if datetime.IsWeekend(now) {
		log.Println("Today is the weekend. Nothing to do.")
		os.Exit(0)
	}

	if datetime.IsHoliday(now) {
		log.Println("Today is a holiday. Nothing to do.")
		os.Exit(0)
	}
}

func initDatabase() (*sql.DB, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	err = database.WaitForDB(db, 10, time.Minute)
	if err != nil {
		database.Close(db)
		return nil, err
	}

	db, err = database.Migrate(db)
	if err != nil {
		database.Close(db)
		return nil, err
	}

	return db, nil
}

func runInitialHistory(db *sql.DB, historyConfig string) {
	currencies := strings.Split(historyConfig, ",")
	for _, currency := range currencies {
		count, countErr := database.CountCurrencyHistory(db, string(currency))
		if countErr != nil {
			log.Printf("Unable to get currency %s count registries", string(currency))
			continue
		}
		if count == 0 {
			data, err := service.GetExchangeHistory(string(currency))
			if err != nil {
				log.Printf("Unable to get currency %s history", string(currency))
				continue
			}

			database.InsertExchangeData(db, "exchange_hist", data)
		}
	}
}

func runManualHistory(db *sql.DB, historyConfig string) {
	currencies := strings.Split(historyConfig, ",")
	for _, currency := range currencies {
		data, err := service.GetExchangeHistory(string(currency))
		if err != nil {
			log.Printf("Unable to get currency %s history", string(currency))
			continue
		}
		database.InsertExchangeData(db, "exchange_hist", data)
	}
}

func runExchangeTracker(db *sql.DB, timeLocation *time.Location, currencyList config.CurrencyConfig) {
	if datetime.IsBeforeHour(time.Now().In(timeLocation), 9) {
		log.Println("Clearing exchange rates.")
		database.ClearExchangeRates(db)
	}

	for {
		if datetime.IsAfterHour(time.Now().In(timeLocation), 9) {
			break
		}
		log.Println("Waiting to start tracking.")
		time.Sleep(5 * time.Minute)
	}

	for {
		if !(datetime.IsBetween(time.Now().In(timeLocation))) {
			log.Println("Finished collecting exchange data")
			break
		}

		data, err := service.GetExchangeValues(currencyList.RealTime)
		if err != nil {
			log.Println(err)
		}

		log.Println(data)

		err = database.InsertExchangeData(db, "exchange_rates", data)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(5 * time.Minute)
	}

	for {
		if datetime.IsAfterHour(time.Now().In(timeLocation), 20) {
			break
		}
		log.Printf("End-of-day collection pending")
		time.Sleep(10 * time.Minute)
	}

	data, err := service.GetExchangesDayValue(currencyList.History)
	if err != nil {
		log.Println(err)
	}

	err = database.InsertExchangeData(db, "exchange_hist", data)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	checkNonWorkingDay()

	timeLocation, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal("Error loading location:", err)
	}

	currencyList := config.GetExchangeList()

	db, err := initDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close(db)

	runHistory := flag.Bool("run-history", false, "Executes only the history data fetch")
	if !flag.Parsed() {
		flag.Parse()
	}
	if *runHistory {
		runManualHistory(db, currencyList.History)
		log.Println("History fetch completed via flag. Exiting.")
		return
	}

	runInitialHistory(db, currencyList.History)

	runExchangeTracker(db, timeLocation, currencyList)
}
