package ethsigdb

import (
	_ "embed"
	"encoding/json"
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
	var dataset map[string]string

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

func (e *ETHSigDB) AddEntries(entries []Entry) {
	for _, entry := range entries {
		e.dataset[entry.TopicHash] = entry.Event
	}
}
