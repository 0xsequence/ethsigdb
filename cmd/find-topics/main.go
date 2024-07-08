package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/0xsequence/ethsigdb"
)

// dedupelist.com

func main() {
	db, _ := ethsigdb.New(nil)
	client, _ := ethsigdb.NewRemoteLookup()
	newEvents := []string{}

	fileData, err := os.ReadFile("./more_topics.csv")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(fileData), "\n")

	for _, topic := range lines {
		_, ok := db.Lookup(strings.TrimSpace(topic))
		if ok {
			continue // we have it in the db, move on
		}

		resp, err := client.FindEventSig(context.Background(), topic)
		if err != nil {
			log.Fatal(fmt.Errorf("topicHash lookup failed for %s: %w", topic, err))
		}

		result := resp.Result.Event

		if len(result) == 0 || len(result[topic]) == 0 {
			continue
		}

		err = db.AddEntries([]ethsigdb.Entry{{
			Event: result[topic][0].Name,
		}})
		if err != nil {
			log.Fatal(fmt.Errorf("add entry failed for %s: %w", topic, err))
		}

		// lets confirm it looks up now..
		_, ok = db.Lookup(strings.TrimSpace(topic))
		if !ok {
			log.Fatal(fmt.Errorf("topicHash lookup failed for %s after new entry of '%s'", topic, result[topic][0].Name))
		}

		newEvents = append(newEvents, result[topic][0].Name)
	}

	// NOTE: the idea is to copy these and past them into ./cmd/build-assets/more_events.csv
	// and rebuild the build-assets which will in turn rebuild ./ethsigdb.json
	for _, ev := range newEvents {
		fmt.Println(ev)
	}
}
