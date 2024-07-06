package main

import (
	"fmt"
	"log"

	"github.com/0xsequence/ethsigdb"
)

func main() {
	sigdb := ethsigdb.Default()

	// Lookup event signature
	topicHash := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	event, ok := sigdb.Lookup(topicHash)
	if !ok {
		log.Fatal("event not found")
	}

	// this will return Transfer(address,address,uint256)
	fmt.Println("found event!", event)
}
