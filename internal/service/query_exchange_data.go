package service

import (
	"ExchangeTrack/internal/config"
	"ExchangeTrack/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func TruncateDate(timestamp string) string {
	ts, _ := strconv.ParseInt(timestamp, 10, 64)

	t := time.Unix(ts, 0)
	dateStr := t.Format("2006-01-02")

	return dateStr
}

func GetExchangeHistory(currency string) ([]model.CurrencyData, error) {
	url := fmt.Sprintf("https://economia.awesomeapi.com.br/json/daily/%s/90", currency)
	var currencyValues []model.CurrencyData
	var data []map[string]string

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

	code := fmt.Sprintf("%s-%s", data[0]["code"], data[0]["codein"])

	for _, entry := range data {
		high, _ := strconv.ParseFloat(entry["high"], 64)
		low, _ := strconv.ParseFloat(entry["low"], 64)
		bid, _ := strconv.ParseFloat(entry["bid"], 64)

		average := (high + low) / 2

		createDate := TruncateDate(entry["timestamp"])

		newCurrencyValue := model.CurrencyData{
			Code:       code,
			Timestamp:  entry["timestamp"],
			CreateDate: createDate,
			Bid:        bid,
			High:       high,
			Low:        low,
			Average:    average,
		}

		currencyValues = append(currencyValues, newCurrencyValue)

	}

	return currencyValues, err

}

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

func GetExchangesDayValue(currencyList string) ([]model.CurrencyData, error) {
	currencyValues, err := GetExchangeValues(currencyList)

	if err != nil {
		log.Println("Unable to get exchange end of day values.")
		return nil, err
	}

	for i := range currencyValues {
		createDate := TruncateDate(currencyValues[i].Timestamp)
		currencyValues[i].CreateDate = createDate
	}

	return currencyValues, nil
}
