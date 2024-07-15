package main

import (
	"context"
	"fmt"

	"github.com/cedro-finance/metalayer-sequencer/relay"
	_ "github.com/cedro-finance/metalayer-sequencer/client"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		fmt.Println("cancelled")
		cancel()
	}()
	go relay.Start(ctx)
	<- ctx.Done()
}
