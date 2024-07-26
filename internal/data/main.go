package data

import (
	_ "github.com/lib/pq"
	_ "gitlab.com/distributed_lab/kit/pgdb"
)

type MasterQ interface {
	New() MasterQ

	Link() TransactionQ

	Transaction(fn func(db MasterQ) error) error
}
