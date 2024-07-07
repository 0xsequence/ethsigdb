package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/0xsequence/ethsigdb"
)

// run with:
// go run . > ../../ethsigdb.json

func main() {
	// topic0 signatures directory
	dir := "/Users/peter/Dev/other/topic0/signatures"

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	db, _ := ethsigdb.New(nil)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v\n", filePath, err)
			continue
		}

		// build a map of file name to content
		// event topic hash :: event definition

		// topicHash := file.Name()
		// if !strings.HasPrefix(topicHash, "0x") {
		// 	topicHash = fmt.Sprintf("0x%s", topicHash)
		// }

		err = db.AddEntries([]ethsigdb.Entry{{
			Event: strings.TrimSpace(string(content)),
		}})
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, event := range moreEvents {
		err = db.AddEntries([]ethsigdb.Entry{{
			Event: event,
		}})
		if err != nil {
			log.Fatal(err)
		}
	}

	out, err := db.DatasetJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
