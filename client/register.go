package client

import (
	"fmt"

	"github.com/cedro-finance/metalayer-sequencer/relay"
)

func init() {
	fmt.Println("register")
	relay.RegisteredClients["ethereum"] = NewClient
	relay.RegisteredClients["bnb"] = NewClient
}