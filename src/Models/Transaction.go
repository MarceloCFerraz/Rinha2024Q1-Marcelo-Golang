package models

import (
	"math"
	"time"
)

type Transaction struct {
	ClientId    *int       `json:"id_cliente"` // optional
	Description string     `json:"descricao"`
	Type        string     `json:"tipo"` // 8 bits long => only ascii chars. Should use rune if need to handle more characters
	Value       float64    `json:"valor"`
	When        *time.Time `json:"realizada_em"` // optional
}

func (t *Transaction) IsInvalid() bool {
	return ((string(t.Type) != "c" && string(t.Type) != "d") ||
		len(t.Type) > 1 ||
		math.Mod(t.Value, 1) != 0 ||
		t.Value <= 0 ||
		len(t.Description) > 10 ||
		len(t.Description) < 1)
}

func (t *Transaction) GetDbOperation() string {
	response := "debit"
	if string(t.Type) == "c" {
		response = "credit"
	}
	return response
}
