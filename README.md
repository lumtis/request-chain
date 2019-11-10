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

### Others

Account used:
pierre: cosmos1en2n7qy77su9c8uw8hgecvfg8czdrfpu2x5qxl
perrine: cosmos1vfpdz6wpyq7f6qw843t2qehqmf42f2fuc7u9n3
