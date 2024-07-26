package data

type TransactionQ interface{}

type CoupleLinks struct {
	Shortened string `db:"shortened" structs:"shortened"`
	Original  string `db:"original" structs:"original"`
}
