package service

import (
	"ExchangeTrack/internal/config"
	"ExchangeTrack/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// func GetExchangeHistory(currency string) {
// 	url := fmt.Sprintf("https://economia.awesomeapi.com.br/json/daily/%s/90", currency)
// 	var currencyValues []model.CurrencyData
// 	var data []map[string]string

// }

func GetExchangeValues(currencyList string) ([]model.CurrencyData, error) {
	url := fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%s", currencyList)
	var currencyValues []model.CurrencyData
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

		newCurrencyValue := model.CurrencyData{
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
