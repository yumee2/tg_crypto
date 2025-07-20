package main

import (
	"fmt"
	"tg-crypto-tracker/internal/adapters"
	"tg-crypto-tracker/internal/infrastructure/parser"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	var klinesData []parser.KlineDataWrapper

	tickers, err := parser.GetAllTickers()
	if err != nil {
		fmt.Println(err)
		return
	}

	var tickerChannel chan parser.KlineDataWrapper = make(chan parser.KlineDataWrapper, len(tickers))

	go parser.ParseTokens(tickers, 100, tickerChannel)

	for data := range tickerChannel {
		klinesData = append(klinesData, data)
	}
	fmt.Println("Ttokens:", len(klinesData))
	fmt.Println("Amount of tickers: ", len(tickers))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://tgcrypto.netlify.app"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.POST("/auth", adapters.AuthUser)
	r.Run()
}
