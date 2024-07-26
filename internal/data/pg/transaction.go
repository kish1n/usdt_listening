package pg

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/kish1n/usdt_listening/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const tableName = "transfers"

func newTransactionQ(db *pgdb.DB) data.TransactionQ {
	return &TransactionQ{
		db:  db,
		sql: sq.StatementBuilder,
	}
}

type TransactionQ struct {
	db  *pgdb.DB
	sql sq.StatementBuilderType
}
