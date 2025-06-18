package main

import (
	"context"
	"main/ticker"
	"main/ticker/core"
	"main/ticker/provider/coingecko"
	"main/ticker/provider/noaa"
	"main/ticker/provider/yahoofinance"
	"main/ticker/provider/youtube"
	"time"
)

func main() {
	rt := ticker.Runtime{
		ConfigurationPath: "./config.json",
		Providers: []core.Provider{
			&core.IntervalProvider{
				Parent:         &noaa.NOAA{},
				UpdateInterval: time.Minute,
			},
			&core.IntervalProvider{
				Parent:         &youtube.YouTube{},
				UpdateInterval: time.Hour,
			},
			&core.IntervalProvider{
				Parent:         &yahoofinance.YahooFinance{},
				UpdateInterval: time.Minute * 2,
			},
			&core.IntervalProvider{
				Parent:         &coingecko.CoinGecko{},
				UpdateInterval: time.Minute * 2,
			},
		},
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
