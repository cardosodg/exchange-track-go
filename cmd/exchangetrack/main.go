package main

import (
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

	db := database.Connect()
	defer database.Close(db)
	database.CreateTables(db)

	for i := 0; i < 3; i++ {
		if !(datetime.IsBetween(time.Now())) {
			log.Println("Finished collecting exchange data")
			break
		}

		data, err := service.GetExchangeValues()
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

	data, err := service.GetExchangeValues()
	if err != nil {
		log.Println(err)
	}

	database.InsertExchangeData(db, "exchange_hist", data)
	if err != nil {
		log.Println(err)
	}

}
