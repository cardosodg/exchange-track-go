package main

import (
	"ExchangeTrack/internal/config"
	"ExchangeTrack/internal/database"
	"ExchangeTrack/internal/datetime"
	"ExchangeTrack/internal/service"
	"fmt"
	"log"
	"os"
	"strings"
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

	database.ClearExchangeRates(db)

	currencies := strings.Split(currencyList.History, ",")
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

	time.Sleep(30 * time.Second)
	data, err := service.GetExchangesDayValue(currencyList.History)
	if err != nil {
		log.Println(err)
	}

	database.InsertExchangeData(db, "exchange_hist", data)
	if err != nil {
		log.Println(err)
	}

}
