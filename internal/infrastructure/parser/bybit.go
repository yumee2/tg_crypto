package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Ticker struct {
	Symbol    string  `json:"symbol"`
	LastPrice float64 `json:"lastPrice,string"`
	HighPrice float64 `json:"highPrice24h,string"`
	LowPrice  float64 `json:"lowPrice24h,string"`
	Volume    float64 `json:"volume24h,string"`
}

type KlineResponse struct {
	Category string     `json:"category"`
	Symbol   string     `json:"symbol"`
	List     [][]string `json:"list"`
}

type KlineData struct {
	StartTime  string
	OpenPrice  string
	HighPrice  string
	LowPrice   string
	ClosePrice string
	Volume     string
	Turnover   string
}

type KlineDataWrapper struct {
	data   []KlineData
	symbol string
}

func GetAllTickers() (map[string]Ticker, error) {
	tickers := make(map[string]Ticker)

	resp, err := http.Get("https://api.bybit.com/v5/market/tickers?category=linear")
	if err != nil {
		return nil, errors.New("cannot get tickers from bybit")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("cannot parse body json from Bybit")
	}

	var APIResponse struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			List []Ticker `json:"list"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &APIResponse); err != nil {
		return nil, errors.New("cannot parse body json from Bybit")
	}

	for _, ticker := range APIResponse.Result.List {
		tickers[ticker.Symbol] = ticker
	}

	return tickers, nil
}

func ParseTokens(tickers map[string]Ticker, out chan<- KlineDataWrapper) {
	var wg sync.WaitGroup
	limiter := time.Tick(1 * time.Second / 100)

	for key := range tickers {
		<-limiter
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			getTokenData(key, 100, out)
		}(key)
	}

	wg.Wait()
	close(out)
}

func getTokenData(symbol string, limit int, out chan<- KlineDataWrapper) {
	url := fmt.Sprintf("https://api.bybit.com/v5/market/kline?category=linear&symbol=%s&interval=5&limit=%d", symbol, limit)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var apiResponse struct {
		Result KlineResponse `json:"result"`
		RetMsg string        `json:"retMsg"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		panic(err)
	}

	var klines []KlineData
	for _, item := range apiResponse.Result.List {
		if len(item) < 7 {
			continue
		}
		klines = append(klines, KlineData{
			StartTime:  item[0],
			OpenPrice:  item[1],
			HighPrice:  item[2],
			LowPrice:   item[3],
			ClosePrice: item[4],
			Volume:     item[5],
			Turnover:   item[6],
		})
	}

	out <- KlineDataWrapper{data: klines, symbol: apiResponse.Result.Symbol}
}
