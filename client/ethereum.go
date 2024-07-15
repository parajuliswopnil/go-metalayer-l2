package client

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/cedro-finance/metalayer-sequencer/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client     *ethclient.Client
	nextHeight uint64
	blockTime  time.Duration
}


func NewClient(rpcUrl common.Url, startHeight uint64) (common.IClient) {
	client := &Client{}

	ethClient, err := ethclient.Dial(string(rpcUrl))
	if err != nil {
		panic(err)
	}

	client.client = ethClient

	return client
}

// queries for the appropriate packets that arise from L1. The appropriate packets
// includes deposit and forced inclusion transactions on L1 bridge contracts
func (c *Client) LoopIteration(ctx context.Context, msgChan chan<- *common.Message) {
	fmt.Println("reached in loop iteration")
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		latestHeight := c.LatestHeight(ctx)
		if latestHeight <= c.nextHeight {
			difference := c.nextHeight - latestHeight + 1 // to handle the equal case
			time.Sleep(time.Second * c.blockTime * time.Duration(difference))
			continue
		}
		msgs, err := c.FilterLogs(ctx, c.nextHeight, latestHeight)
		if err != nil {
			// warn
			continue
		}
		for _, msg := range msgs {
			msgChan <- msg
		}
		c.nextHeight = latestHeight + 1 // filter logs in inclusive of start and end height
	}
}

func (c *Client) FilterLogs(ctx context.Context, from, to uint64) ([]*common.Message, error) {
	return []*common.Message{
		{From: ethCommon.HexToAddress(""), To: ethCommon.HexToAddress(""), Value: big.NewInt(100)},
		{From: ethCommon.HexToAddress(""), To: ethCommon.HexToAddress(""), Value: big.NewInt(200)},
	}, nil
}

func (c *Client) LatestHeight(ctx context.Context) uint64 {
	latestHeight, err := c.client.BlockNumber(ctx)
	if err != nil {
		return 0
	}
	return latestHeight
}
