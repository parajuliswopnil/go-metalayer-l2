package relay

import (
	"context"
	"fmt"
	"time"

	"github.com/cedro-finance/metalayer-sequencer/common"
)

var (
	RegisteredClients = make(map[string]func(common.Url, uint64) common.IClient)
	ChainClients      = make(map[string]common.IClient)
)

const (
	BatchSize = 10
	SequencingWindow = time.Minute
)

func Start(ctx context.Context) {
	msgCh := make(chan *common.Message)
	for k, v := range RegisteredClients {
		ChainClients[k] = v(common.Url("https://eth-sepolia.g.alchemy.com/v2/GLYjGj3YzX6ayV_RlGe--fhdxbcuyDBv"), 6315600)
	}

	go processMessages(ctx, msgCh)
	for _, v := range ChainClients {
		go v.LoopIteration(ctx, msgCh)
	}
}

func processMessages(ctx context.Context, msgChan <-chan *common.Message) {
	fmt.Println("reached in processor")
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgChan:
			fmt.Println(msg.Value)
		default:
			time.Sleep(time.Second)
		}
	}
}
