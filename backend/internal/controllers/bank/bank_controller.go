package bank

import "log/slog"

type BankController struct {
	Logger *slog.Logger
}

func NewBankController(logger *slog.Logger) *BankController {
	return &BankController{
		Logger: logger,
	}
}

