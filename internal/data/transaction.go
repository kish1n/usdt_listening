package data

import "time"

type TransactionQ interface {
	Insert(trn TransactionData) (*TransactionData, error)
}

type TransactionData struct {
	From_address string    `db:"from_address"`
	To_address   string    `db:"to_address"`
	Value        int64     `db:"value"`
	Id           string    `db:"id"`
	Timestamp    time.Time `db:"timestamp"`
}
