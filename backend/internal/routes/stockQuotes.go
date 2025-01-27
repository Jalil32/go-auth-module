package routes

import (
	"fmt"

	"github.com/piquette/finance-go/quote"
)

func GetStockQuote(symbol string) (float64, error) {
	q, err := quote.Get(symbol)

	if err != nil {
		return -1, fmt.Errorf("Failed to get quote: %v", err)
	}
	if q == nil {
		return -1, fmt.Errorf("No quote found for symbol: %v", symbol)
	}

	return q.Ask, nil
}