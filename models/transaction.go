package models

type Transaction struct {
	Recipient int     `json:"recipient"`
	Sender    int     `json:"sender"`
	Amount    float64 `json:"amount"`
}
