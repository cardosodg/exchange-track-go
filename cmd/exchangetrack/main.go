package main

import (
	"ExchangeTrack/internal/config"
	"ExchangeTrack/internal/database"
	"ExchangeTrack/internal/datetime"
	"ExchangeTrack/internal/service"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if datetime.IsWeekend(time.Now()) {
		log.Println("Today is the weekend. Nothing to do.")
		os.Exit(0)
	}

	if datetime.IsHoliday(time.Now()) {
		log.Println("Today is a holiday. Nothing to do.")
		os.Exit(0)
	}

	currencyList := config.GetExchangeList()

	db := database.Connect()
	defer database.Close(db)
	database.CreateTables(db)

	for {
		if !(datetime.IsBetween(time.Now())) {
			log.Println("Finished collecting exchange data")
			break
		}

		data, err := service.GetExchangeValues(currencyList.RealTime)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(data)

		err = database.InsertExchangeData(db, "exchange_rates", data)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(5 * time.Minute)
	}

	data, err := service.GetExchangeValues(currencyList.History)
	if err != nil {
		log.Println(err)
	}

	database.InsertExchangeData(db, "exchange_hist", data)
	if err != nil {
		log.Println(err)
	}

}
