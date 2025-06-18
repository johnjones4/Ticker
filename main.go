package main

import (
	"context"
	"log/slog"
	"main/ticker"
	"main/ticker/core"
	"main/ticker/provider/coingecko"
	"main/ticker/provider/noaa"
	"main/ticker/provider/yahoofinance"
	"main/ticker/provider/youtube"
	"time"
)

func main() {
	log := slog.Default()

	rt := ticker.Runtime{
		Log:               log,
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
		log.Error("error during init", slog.Any("error", err))
		return
	}

	err = rt.Start(ctx)
	if err != nil {
		log.Error("error during runtime", slog.Any("error", err))
		return
	}
}
