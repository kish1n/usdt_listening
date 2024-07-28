package handlers

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/structs"
	"github.com/google/jsonapi"
	"github.com/kish1n/usdt_listening/internal/data"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
	"gitlab.com/distributed_lab/ape"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

type Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Index int
	Time  time.Time
}

var Client *ethclient.Client

func InitEthereumClient(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)

	ProjectID := os.Getenv("API_KEY")
	logger.Info(ProjectID)
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + ProjectID)

	if err != nil {
		logger.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	Client = client
	logger.Infof("Connected to Ethereum client")
}

func ListenForTransfers(w http.ResponseWriter, r *http.Request) {
	logger := helpers.Log(r)

	ProjectID := os.Getenv("API_KEY")
	logger.Info(ProjectID)
	client, err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/" + ProjectID)

	if err != nil {
		logger.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	Client = client
	logger.Infof("Connected to Ethereum client")

	contractAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)

	sub, err := Client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		logger.Fatalf("Failed to subscribe to logs: %v", err)
	}

	contractABIJSON, err := helpers.ReadABIFile("contractABI.json")
	if err != nil {
		logger.Fatalf("Failed to read contract ABI file: %v", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		logger.Fatalf("Failed to parse contract ABI: %v", err)
	}

	db := helpers.DB(r)

	for {
		select {
		case err := <-sub.Err():
			logger.Fatalf("Error: %v", err)
		case vLog := <-logs:
			logger.Infof("Log: %v", vLog)

			var transferEvent Transfer
			err := contractABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				logger.Fatalf("Failed to unpack log: %v", err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
			logger.Infof("success unpack log %s to %s", transferEvent, transferEvent.To)
			stmt := data.TransactionData{
				FromAddress: transferEvent.From.Hex(),
				ToAddress:   transferEvent.To.Hex(),
				Value:       transferEvent.Value.Int64(),
				Id:          helpers.GenerateUUID(),
				Timestamp:   time.Now(),
			}

			test := structs.Map(stmt)
			logger.Infof("test %s", test)

			res, err := db.Link().Insert(stmt)

			if err != nil {
				logger.WithError(err).Debug("Server error")
				ape.Render(w, &jsonapi.ErrorObject{
					Status: "500",
					Title:  "Server error 500",
					Detail: "Unpredictable behavior",
				})
				return
			}

			logger.Infof("Transfer event: from %s to %s for %d tokens", res.FromAddress,
				res.ToAddress, res.Value)
		}
	}
}
