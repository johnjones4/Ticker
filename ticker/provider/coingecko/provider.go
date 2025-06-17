package coingecko

import (
	"context"
	"fmt"
	"main/ticker/core"
	"strings"
)

type CoinGecko struct {
	coins        []string
	currencyCode string
	client       ClientWithResponsesInterface
	coinNameMap  map[string]string
}

func (p *CoinGecko) Init(ctx context.Context, cfg *core.Configuration) error {
	p.coins = cfg.CryptoCoins
	p.currencyCode = cfg.CurrencyCode
	var err error
	p.client, err = NewClientWithResponses("https://api.coingecko.com/api/v3")
	if err != nil {
		return err
	}
	err = p.loadCoinNames(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *CoinGecko) Update(ctx context.Context) ([]string, error) {
	out := make([]string, 0, len(p.coins))
	ids := strings.Join(p.coins, ",")
	t := true
	res, err := p.client.SimplePriceWithResponse(ctx, &SimplePriceParams{
		Ids:               &ids,
		VsCurrencies:      p.currencyCode,
		Include24hrChange: &t,
	})
	if err != nil {
		return nil, err
	}

	if res.JSON200 == nil {
		return nil, nil
	}

	for _, coin := range p.coins {
		value, ok := (*res.JSON200)[coin]
		if !ok || value.Usd24hChange == nil {
			continue
		}
		var direction string
		if *value.Usd24hChange > 0 {
			direction = "+"
		}
		message := fmt.Sprintf("%s: %s$%0.2f", p.coinNameMap[coin], direction, *value.Usd24hChange)
		out = append(out, message)
	}

	return out, nil
}

func (p *CoinGecko) loadCoinNames(ctx context.Context) error {
	res, err := p.client.CoinsListWithResponse(ctx, &CoinsListParams{})
	if err != nil {
		return err
	}
	if res.JSON200 == nil {
		return nil
	}
	p.coinNameMap = make(map[string]string)
	for _, coin := range *res.JSON200 {
		p.coinNameMap[*coin.Id] = *coin.Name
	}
	return nil
}
