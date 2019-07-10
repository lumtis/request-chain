![Request Chain](requestChain.png)

# Request Chain

Proof of Stake blockchain using Request Network protocol.



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
