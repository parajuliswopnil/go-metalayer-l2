package client

import (
	"context"
	"testing"

	"github.com/cedro-finance/metalayer-sequencer/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

const (
	url = "https://eth-mainnet.g.alchemy.com/v2/GLYjGj3YzX6ayV_RlGe--fhdxbcuyDBv"
)

func TestNewClient(t *testing.T) {
	startHeight := 10
	client := NewClient(common.Url(url), uint64(startHeight))

	blockNumber := client.(*Client).LatestHeight(context.Background())
	assert.NotZero(t, blockNumber)
}

func TestLoopIteration(t *testing.T) {
	cl, err := ethclient.Dial(url)
	assert.NoError(t, err)
	client := &Client{
		client: cl,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	latestHeight := client.LatestHeight(ctx)
	assert.NotEqual(t, uint64(0), latestHeight)

	client.nextHeight = latestHeight - 10

	msgChan := make(chan *common.Message)

	go client.LoopIteration(ctx, msgChan)

	for range 2 {
		msg := <-msgChan
		println(msg.Value)
	}
}
