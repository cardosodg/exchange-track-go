package service

import (
	"ExchangeTrack/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// type ExchangeData struct {
// 	USDBRL ExchangeRate `json:"USDBRL"`
// 	EURBRL ExchangeRate `json:"EURBRL"`
// }

// type ExchangeRate struct {
// 	Bid        json.Number `json:"bid"`
// 	High       json.Number `json:"high"`
// 	Low        json.Number `json:"low"`
// 	Timestamp  string      `json:"timestamp"`
// 	CreateDate string      `json:"create_date"`
// }

type CurrencyData struct {
	Code       string
	Timestamp  string
	CreateDate string
	Bid        float64
	High       float64
	Low        float64
	Average    float64
}

func check_time() {
	start := time.Now()

	time.Sleep(2 * time.Second)

	end := time.Now()

	delta := end.Sub(start)
	fmt.Printf("Time delta: %v\n", delta)
}

func Initialize() {
	currentTime := time.Now()
	fmt.Println("Current Time in String: ", currentTime.Format("2006-01-02T15:04:05"))
	fmt.Println("ISO8601: ", time.Now().Format(time.RFC3339))

	time.Sleep(5 * time.Second)
	check_time()

	// date := time.Date(2025, 3, 18, 0, 0, 0, 0, time.UTC)
	// holidays, err := datetime.GetHolidays("2025")
	// checkHoliday := datetime.IsHoliday(date, holidays)

	// if err != nil {
	// 	fmt.Println("Erro:", err)
	// 	return
	// }

	// if checkHoliday {
	// 	fmt.Printf("A data %s é um feriado.\n", date.Format("2006-01-02"))
	// } else {
	// 	fmt.Printf("A data %s não é um feriado.\n", date.Format("2006-01-02"))
	// }
}

func GetExchangeValues() ([]CurrencyData, error) {
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL,EUR-BRL"
	var currencyValues []CurrencyData
	var data map[string]map[string]string

	client := &http.Client{}
	apiKey := config.GetApiKey()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, entry := range data {
		code := fmt.Sprintf("%s-%s", entry["code"], entry["codein"])
		high, err := strconv.ParseFloat(entry["high"], 64)

		if err != nil {
			log.Printf("Skipping %s: invalid high value", code)
			continue
		}

		low, err := strconv.ParseFloat(entry["low"], 64)

		if err != nil {
			log.Printf("Skipping %s: invalid low value", code)
			continue
		}

		bid, err := strconv.ParseFloat(entry["bid"], 64)
		if err != nil {
			log.Printf("Skipping %s: invalid bid value", code)
			continue
		}

		average := (high + low) / 2

		newCurrencyValue := CurrencyData{
			Code:       code,
			Timestamp:  entry["timestamp"],
			CreateDate: entry["create_date"],
			Bid:        bid,
			High:       high,
			Low:        low,
			Average:    average,
		}

		currencyValues = append(currencyValues, newCurrencyValue)

	}

	return currencyValues, err
}

func ExecuteReadings() {
	var data []CurrencyData

	for i := range 2 {
		fmt.Printf("Executing iteration %d\n", i+1)
		data, _ = GetExchangeValues()

		for _, item := range data {
			currency, err := json.MarshalIndent(item, "", " ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(currency))
		}
		fmt.Println(" ")
		time.Sleep(5 * time.Second)
	}
}

func FinishReadings() {
	fmt.Printf("Waiting 10 seconds to finish execution.\n")
	time.Sleep(10 * time.Second)
	var data []CurrencyData

	data, _ = GetExchangeValues()

	for _, item := range data {
		currency, err := json.MarshalIndent(item, "", " ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(currency))
	}
	fmt.Println("End.")
}
