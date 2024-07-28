package data

import "time"

type TransactionQ interface {
	Insert(trn TransactionData) (*TransactionData, error)
	SortByParameter(address string, parameter string) ([]TransactionData, error)
}

type TransactionData struct {
	FromAddress string    `db:"from_address" json:"id"`
	ToAddress   string    `db:"to_address"   json:"from_address"`
	Value       int64     `db:"value"        json:"to_address"`
	Id          string    `db:"id"           json:"value"`
	Timestamp   time.Time `db:"timestamp"    json:"timestamp"`
}
