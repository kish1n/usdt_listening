package data

import (
	_ "github.com/lib/pq"
	_ "gitlab.com/distributed_lab/kit/pgdb"
)

type MasterQ interface {
	NewMaster() MasterQ

	NewTransaction() TransactionQ

	Transaction(fn func(db MasterQ) error) error
}
