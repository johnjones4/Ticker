package core

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	MarketSymbols []string `json:"marketSymbols"`
	CryptoCoins   []string `json:"cryptoCoins"`
	CurrencyCode  string   `json:"currencyCode"`
}

func (cfg *Configuration) Load(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return err
	}

	return nil
}
