package main

import (
	"context"
	"fmt"
	"log"

	"github.com/0xsequence/ethsigdb"
)

func main() {
	db, _ := ethsigdb.New(nil)
	client, _ := ethsigdb.NewRemoteLookup()
	newEvents := []string{}

	for _, topic := range moreTopics {
		_, ok := db.Lookup(topic)
		if ok {
			continue // we have it in the db, move on
		}

		resp, err := client.FindEventSig(context.Background(), topic)
		if err != nil {
			log.Fatal(err)
		}

		result := resp.Result.Event

		if len(result) == 0 || len(result[topic]) == 0 {
			continue
		}

		err = db.AddEntries([]ethsigdb.Entry{{
			Event: result[topic][0].Name,
		}})
		if err != nil {
			log.Fatal(err)
		}
		newEvents = append(newEvents, result[topic][0].Name)
	}

	// NOTE: the idea is to copy these and past them into ./cmd/build-assets/more_events.go
	// and rebuild the build-assets which will in turn rebuild ./ethsigdb.json
	for _, ev := range newEvents {
		fmt.Printf("\"%s\",\n", ev)
	}
}
