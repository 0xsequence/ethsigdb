package ethsigdb

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/0xsequence/ethkit/ethcoder"
)

//go:embed ethsigdb.json
var embeddedJSON []byte

type ETHSigDB struct {
	dataset map[string]string // event topic hash :: event definition
}

type Entry struct {
	TopicHash string `json:"topicHash"`
	Event     string `json:"event"`
}

func New(jsonDataset []byte) (*ETHSigDB, error) {
	dataset := map[string]string{}

	if len(jsonDataset) > 0 {
		err := json.Unmarshal(jsonDataset, &dataset)
		if err != nil {
			return nil, err
		}
	}

	return &ETHSigDB{
		dataset: dataset,
	}, nil
}

func Default() *ETHSigDB {
	db, _ := New(embeddedJSON)
	return db
}

func (e *ETHSigDB) Lookup(topicHash string) (string, bool) {
	event, ok := e.dataset[topicHash]
	return event, ok
}

func (e *ETHSigDB) AddEntries(entries []Entry) error {
	for _, entry := range entries {
		topicHash, eventSig, err := ethcoder.EventTopicHash(entry.Event)
		if err != nil {
			return fmt.Errorf("invalid entry %s: %w", entry.Event, err)
		}
		e.dataset[topicHash.String()] = eventSig
	}

	return nil
}

func (e *ETHSigDB) WriteToFile(filepath string) error {
	out, err := e.DatasetJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, out, 0644)
}

func (e *ETHSigDB) DatasetJSON() ([]byte, error) {
	return json.MarshalIndent(e.dataset, "", "  ")
}
