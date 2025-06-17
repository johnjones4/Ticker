package yahoofinance

import (
	"context"
	"fmt"
	"main/ticker/core"
	"net/http"
)

type YahooFinance struct {
	symbols []string
	client  ClientWithResponsesInterface
}

func (p *YahooFinance) Init(ctx context.Context, cfg *core.Configuration) error {
	p.symbols = cfg.MarketSymbols
	var err error
	p.client, err = NewClientWithResponses("https://query1.finance.yahoo.com", WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Referer", "https://finance.yahoo.com")
		return nil
	}))
	if err != nil {
		return err
	}
	return nil
}

func (p *YahooFinance) Update(ctx context.Context) ([]string, error) {
	out := make([]string, 0, len(p.symbols))
	for _, symbol := range p.symbols {
		res, err := p.client.GetV8FinanceChartSymbolWithResponse(ctx, symbol, &GetV8FinanceChartSymbolParams{
			Range:    GetV8FinanceChartSymbolParamsRangeN1d,
			Interval: GetV8FinanceChartSymbolParamsIntervalN30m,
		})
		if err != nil {
			return nil, err
		}

		if res.JSON200 == nil || res.JSON200.Chart == nil || len(*res.JSON200.Chart.Result) == 0 || (*res.JSON200.Chart.Result)[0].Meta == nil || (*res.JSON200.Chart.Result)[0].Meta.ShortName == nil || (*res.JSON200.Chart.Result)[0].Meta.PreviousClose == nil || (*res.JSON200.Chart.Result)[0].Meta.RegularMarketPrice == nil {
			continue
		}

		name := *(*res.JSON200.Chart.Result)[0].Meta.ShortName
		prevClose := *(*res.JSON200.Chart.Result)[0].Meta.PreviousClose
		current := *(*res.JSON200.Chart.Result)[0].Meta.RegularMarketPrice
		pcnt := ((current - prevClose) / prevClose) * 100
		var direction string
		if pcnt > 0 {
			direction = "+"
		}
		message := fmt.Sprintf("%s: %s%0.2f%%", name, direction, pcnt)

		out = append(out, message)
	}
	return out, nil
}
