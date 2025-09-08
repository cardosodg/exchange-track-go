package model

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
