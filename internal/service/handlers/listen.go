package handlers

import (
	"context"
	"github.com/kish1n/usdt_listening/internal/gen"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kish1n/usdt_listening/internal/config"
	"github.com/kish1n/usdt_listening/internal/data"
	"github.com/kish1n/usdt_listening/internal/data/pg"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
)

func ListenTransfers(cfg config.Config) {
	logger := cfg.Log()

	ProjectID := cfg.ServiceConfig().TokenKey

	logger.Info(ProjectID)
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/" + ProjectID)
	if err != nil {
		logger.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return
	}

	logger.Infof("Connected to Ethereum client")

	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		logger.Fatalf("Failed to subscribe to logs: %v", err)
		return
	}

	db := pg.NewMasterQ(cfg.DB())

	for {
		select {
		case err := <-sub.Err():
			logger.Fatalf("Error: %v", err)
		case vLog := <-logs:
			logger.Infof("Log: %v", vLog)

			storageInstance, err := gen.NewStorage(contractAddress, client)
			if err != nil {
				logger.Fatalf("Failed to instantiate contract: %v", err)
				return
			}

			event, err := storageInstance.ParseTransfer(vLog)
			if err != nil {
				logger.Fatalf("Failed to parse transfer event: %v", err)
				return
			}

			stmt := data.TransactionData{
				Sender:    event.From.Hex(),
				Recipient: event.To.Hex(),
				Value:     event.Value.Int64(),
				Id:        helpers.GenerateUUID(),
				Timestamp: time.Now(),
			}

			res, err := db.NewTransaction().Insert(stmt)
			if err != nil {
				logger.WithError(err).Error("Failed to insert transaction data")
				return
			}

			logger.Infof("Transfer event: from %s to %s for %d tokens", res.Sender,
				res.Recipient, res.Value)
		}
	}
}
