package handlers

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kish1n/usdt_listening/internal/service/helpers"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
)

type Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
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

func readABIFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
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

	contractABIJSON, err := readABIFile("/usr/local/bin/contractABI.json")
	if err != nil {
		logger.Fatalf("Failed to read contract ABI file: %v", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(contractABIJSON))
	if err != nil {
		logger.Fatalf("Failed to parse contract ABI: %v", err)
	}

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

			logger.Infof("Transfer event: from %s to %s for %d tokens", transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Value)
		}
	}
}
