package data

import (
	"time"
)

type TransactionQ interface {
	Insert(trn TransactionData) (*TransactionData, error)
	FilterBySender(address string) ([]TransactionData, error)
	FilterByRecipient(address string) ([]TransactionData, error)
}

type TransactionData struct {
	Sender    string    `db:"sender"       json:"sender"`
	Recipient string    `db:"recipient"    json:"recipient"`
	Value     int64     `db:"value"        json:"value"`
	Id        string    `db:"id"           json:"id"`
	Timestamp time.Time `db:"timestamp"    json:"timestamp"`
}
