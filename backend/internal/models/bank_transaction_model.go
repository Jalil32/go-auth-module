package models

import "time"

type Transaction struct {
	ID          string
	Date        time.Time
	Amount      float64
	Description string
}
