package common

import (
	"context"
	"math/big"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

type Url string

type Message struct {
	From  ethCommon.Address
	To    ethCommon.Address
	Value *big.Int
}

type IClient interface {
	LoopIteration(ctx context.Context, msgChan chan<- *Message)
}
