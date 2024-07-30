package pg

import (
	"github.com/kish1n/usdt_listening/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
	"log"
)

func NewMasterQ(db *pgdb.DB) data.MasterQ {
	dataBase := db.Clone()
	log.Println("db clone")
	return &masterQ{
		db: dataBase,
	}
}

type masterQ struct {
	db *pgdb.DB
}

func (m *masterQ) NewMaster() data.MasterQ {
	return NewMasterQ(m.db)
}

func (m *masterQ) NewTransaction() data.TransactionQ {
	return newTransactionQ(m.db)
}

func (m *masterQ) Transaction(fn func(q data.MasterQ) error) error {
	return m.db.Transaction(func() error {
		return fn(m)
	})
}
