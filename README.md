![Request Chain](requestChain.png)

# Request Chain

Example of a Proof of Stake blockchain to be used for the storage layer of the Request protocol

### Message


```
type MsgAppendBlock struct {
	Block string
}
```

### Query

```
type QueryGetBlock struct {
	Hash string
}
```

```
type QueryGetBlockCount struct {
}
```
