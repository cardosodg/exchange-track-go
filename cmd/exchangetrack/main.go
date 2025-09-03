package main

import (
	"ExchangeTrack/internal/service"
)

func main() {
	service.Initialize()
	service.ExecuteReadings()
	service.FinishReadings()

	// // Carregar as configurações do arquivo .env
	// cfg := config.LoadConfig()

	// // Conectar ao banco de dados
	// db := database.Connect(cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName)
	// defer database.Close(db) // Fechar a conexão quando terminar

	// // Criando um novo ExchangeRateDB para inserir
	// rate := database.ExchangeDataDB{
	// 	Bid:        "5.35",
	// 	Timestamp:  "1743556263",
	// 	CreateDate: "2025-04-01 22:11:03",
	// 	Type:       "USDBRL",
	// }

	// // Inserir os dados no banco de dados
	// err := database.InsertExchangeData(db, rate)
	// if err != nil {
	// 	log.Fatal("Erro ao inserir dados:", err)
	// } else {
	// 	fmt.Println("Dado inserido com sucesso!")
	// }

}
