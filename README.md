ethsigdb
========

Ethereum transaction log event signature database and Go lookup library.

To purpose of this dataset is to do a reverse lookup to convert an Ethereum
transaction log topicHash to the event defintion.

For example, when an ERC20 transfer transaction is mined, it will emit the
event log with topics[0] value of: `0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`

In order to look up from `0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef` what is
the actual event definition, we use this library and database to return:
`Transfer(address,address,uint256)`.

This is useful when converting event data to human readable formats, or
for decoding an event.


## Dataset

The signature dataset was originally sourced from https://github.com/otterscan/topic0
which in turn is sourced from https://github.com/ethereum/sourcify. 

Signature data was copied into ./cmd/build-dataset/more_events_1.csv and additional
events are specified in ./cmd/build-dataset/more_events_2.csv. When running the build-dataset
program, it will compute the topic hashes and store all datasets into ./ethsigdb.json.

Finally, if you'd like to search for events by topic hash, you can use ./cmd/find-topics
which will use the remote openchain.xyz service to look up events by their topic hash,
and them print any topic hashes which are not in our local database, so you can copy
them over to ./cmd/build-dataset/more_events_2.csv and rebuild the local db.


## Example

```go
sigdb := ethsigdb.Default()

// Lookup event signature
topicHash := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

event, ok := sigdb.Lookup(topicHash)
if !ok {
  log.Fatal("event not found")
}

// this will return Transfer(address,address,uint256)
fmt.Println("found event!", event)
```


## Caveats

There are a few caveats when using this database to decode generic events.

Event arguments may contain "indexed" arguments which are not necessary in the
order of arguments, for example an event like, `Deposit(address indexed from, uint256 amount, uint256 expiry indexed)`
has two indexed arguments, but the topic hash for this event will be identical to `Deposit(address,uint256,uint256)`,
yet without knowing that the first and third arguments are the ones which are indexed, then its not possible
to do a reverse lookup and decode. For emitted events which have serial indexed arguments, then decoders
can try that path, but this is not a reliable assumption. Therefore, the best is to know the full event
definition which includes the indexed arguments.


## LICENSE

Apache 2.0

