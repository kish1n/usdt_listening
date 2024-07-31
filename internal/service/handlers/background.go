package handlers

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
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
