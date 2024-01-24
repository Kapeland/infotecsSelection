package types

import "time"

type Operation struct {
	Time   time.Time `json:"time"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount float64   `json:"amount"`
}
