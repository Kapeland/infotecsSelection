package types

type Wallet struct {
	Id      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type WltForSend struct {
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}
