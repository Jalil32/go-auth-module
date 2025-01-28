package stock

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piquette/finance-go/quote"
)

type StockController struct {
	Logger *slog.Logger
}

func NewStockController(logger *slog.Logger) *StockController {
	return &StockController{
		Logger: logger,
	}
}

// GetStockQuoteHandler handles the HTTP request and response for fetching a stock quote.
func (s *StockController) GetStockQuoteHandler(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	stockQuote, err := s.GetStockQuote(symbol)
	if err != nil {
		s.Logger.Error("Failed to get stock quote", "symbol", symbol, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"symbol": symbol, "quote": stockQuote})
}

// GetStockQuote fetches the stock quote for a given symbol.
func (s *StockController) GetStockQuote(symbol string) (float64, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		s.Logger.Error("Failed to get stock quote", "symbol", symbol, "error", err)
		return -1, fmt.Errorf("failed to get quote: %v", err)
	}
	if q == nil {
		s.Logger.Error("No quote found for symbol", "symbol", symbol)
		return -1, fmt.Errorf("no quote found for symbol: %v", symbol)
	}

	return q.Ask, nil
}
