package main

import (
	"context"
	"main/ticker"
	"main/ticker/core"
	"main/ticker/output"
	"main/ticker/provider/coingecko"
	"main/ticker/provider/yahoofinance"
	"time"
)

func main() {
	rt := ticker.Runtime{
		ConfigurationPath: "./config.json",
		Providers: []core.Provider{
			&core.IntervalProvider{
				Parent:         &yahoofinance.YahooFinance{},
				UpdateInterval: time.Minute * 2,
			},
			&core.IntervalProvider{
				Parent:         &coingecko.CoinGecko{},
				UpdateInterval: time.Minute * 2,
			},
		},
		Output: &output.LogOutput{},
	}

	ctx := context.Background()
	err := rt.Init(ctx)
	if err != nil {
		panic(err)
	}

	err = rt.Start(ctx)
	if err != nil {
		panic(err)
	}
}
