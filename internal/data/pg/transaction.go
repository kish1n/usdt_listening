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

//func ListenForTransfers(db *sql.DB, contractAddress common.Address, contractABI string) {
//	query := ethereum.FilterQuery{
//		Addresses: []common.Address{contractAddress},
//	}
//
//	logs := make(chan types.Log)
//	sub, err := Client.SubscribeFilterLogs(context.Background(), query, logs)
//	if err != nil {
//		log.Fatalf("Failed to subscribe to logs: %v", err)
//	}
//
//	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
//	if err != nil {
//		log.Fatalf("Failed to parse contract ABI: %v", err)
//	}
//
//	for {
//		select {
//		case err := <-sub.Err():
//			log.Fatalf("Error: %v", err)
//		case vLog := <-logs:
//			var transferEvent Transfer
//			err := parsedABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
//			if err != nil {
//				log.Fatalf("Failed to unpack log: %v", err)
//			}
//
//			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
//			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
//
//			log.Printf("Transfer event: from %s to %s for %d tokens", transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Value)
//
//			_, err = db.Exec(
//				"INSERT INTO transfers (from_address, to_address, value) VALUES ($1, $2, $3)",
//				transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Value,
//			)
//			if err != nil {
//				log.Fatalf("Failed to insert transfer event into database: %v", err)
//			}
//		}
//	}
//}
