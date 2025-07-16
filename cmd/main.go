package main

import (
	"fmt"
	"tg-crypto-tracker/internal/infrastructure/parser"
)

func main() {

	tickers, err := parser.GetAllTickers()
	if err != nil {
		fmt.Println(err)
		return
	}

	var tickerChannel chan parser.KlineDataWrapper = make(chan parser.KlineDataWrapper, len(tickers))

	go parser.ParseTokens(tickers, 100, tickerChannel)

	counter := 0
	for data := range tickerChannel {
		counter++
		fmt.Println("Processing data:", data)
	}
	fmt.Println("Ttokens:", counter)
	fmt.Println("Amount of tickers: ", len(tickers))
	// 	r := gin.Default()
	// 	r.Run()
}
