package pg

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/kish1n/usdt_listening/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
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

func (q *TransactionQ) Insert(trn data.TransactionData) (*data.TransactionData, error) {
	clauses := structs.Map(trn)
	var result data.TransactionData
	stmt := sq.Insert(tableName).SetMap(clauses).Suffix("returning *")
	err := q.db.Get(&result, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert transaction to db")
	}
	return &result, nil
}

func (q *TransactionQ) SortByParameter(address string, parameter string) ([]data.TransactionData, error) {
	var result []data.TransactionData
	stmt := sq.Select("*").From(tableName).Where(sq.Eq{parameter: address})
	err := q.db.Select(&result, stmt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select by origin link in db")
	}

	return result, nil
}
