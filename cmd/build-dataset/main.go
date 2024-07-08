package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/0xsequence/ethsigdb"
)

// run with:
// go run . > ../../ethsigdb.json

func main() {
	files := []string{
		"./more_events_1.csv", "./more_events_2.csv",
	}

	db, _ := ethsigdb.New(nil)

	for _, file := range files {
		fileData, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		lines := strings.Split(string(fileData), "\n")

		for _, event := range lines {
			err = db.AddEntries([]ethsigdb.Entry{{
				Event: strings.TrimSpace(event),
			}})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	out, err := db.DatasetJSON()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
