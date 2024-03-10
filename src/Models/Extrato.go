package models

import "time"

type Report struct {
	Balance          Total         `json:"saldo"`
	LastTransactions []Transaction `json:"ultimas_transacoes"`
}

type Total struct {
	// customer_balance, customer_limit, report_date, last_transactions
	Total      int       `json:"total"`
	ReportDate time.Time `json:"data_extrato"`
	Limit      int       `json:"limite"`
}
