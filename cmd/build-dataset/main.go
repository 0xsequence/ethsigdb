package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// topic0 signatures directory
	dir := "/Users/peter/Dev/other/topic0/signatures"

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	db := map[string]string{}

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

		topicHash := file.Name()
		if !strings.HasPrefix(topicHash, "0x") {
			topicHash = fmt.Sprintf("0x%s", topicHash)
		}

		db[topicHash] = strings.TrimSpace(string(content))
	}

	out, err := json.Marshal(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
