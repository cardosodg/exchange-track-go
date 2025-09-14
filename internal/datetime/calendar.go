package datetime

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Holiday struct {
	Date string `json:"date"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func IsBetween(date time.Time) bool {
	start := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location())
	end := time.Date(date.Year(), date.Month(), date.Day(), 17, 0, 0, 0, date.Location())

	return (date.Equal(start) || date.After(start)) && (date.Equal(end) || date.Before(end))
}

func IsWeekend(date time.Time) bool {
	weekday := date.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func GetHolidays(year int) ([]Holiday, error) {
	// const url = "https://brasilapi.com.br/api/feriados/v1/2025"
	url := fmt.Sprintf("https://brasilapi.com.br/api/feriados/v1/%d", year)
	var holidays []Holiday

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error in getting holidays: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return holidays, nil
}

func InitializeData() {
	// currentYear := fmt.Sprintf("%d", time.Now().Year())
	currentYear := time.Now().Year()
	holidays, err := GetHolidays(currentYear)
	fmt.Println(holidays)
	fmt.Println(err)
}

func IsHoliday(date time.Time) bool {
	dateStr := date.Format("2006-01-02")

	holidays, err := GetHolidays(date.Year())
	if err != nil {
		log.Println(err)
	}

	for _, holiday := range holidays {
		if holiday.Date == dateStr {
			return true

		}
	}
	return false
}
