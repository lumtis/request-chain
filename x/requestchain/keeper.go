package requestchain

import (
	"encoding/json"
	// "fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/list"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// Record the content of a block with the associated timestamp
type BlockRecord struct {
	Block string			`json:"block"`
	Timestamp int64		`json:"timestamp"`
}

// Index of the block with associated timestamp
type BlockIndex struct {
	Index uint64			`json:"index"`
	Timestamp int64		`json:"timestamp"`
}

// NewKeeper creates new instances of the requestchain Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Get a block from the store from its hash
func (k Keeper) GetBlock(ctx sdk.Context, index uint64) []byte {
	// Get the store as a lsit
	store := ctx.KVStore(k.storeKey)
	list := list.NewList(k.cdc, store)
	var blockRecord string

	if list.Len() <= index {
		return []byte{}
	}

	err := list.Get(index, &blockRecord)
	if err != nil {
    panic(err)
  }

	// blockString := fmt.Sprintf("%v", block)

	return []byte(blockRecord)
}

// Get number of blocks
func (k Keeper) GetBlockCount(ctx sdk.Context) uint64 {
	// Get the store as a lsit
	store := ctx.KVStore(k.storeKey)
	list := list.NewList(k.cdc, store)

	return list.Len()
}

// Appends a block of data and returns the hash as an index
func (k Keeper) AppendBlock(ctx sdk.Context, block string) []byte {
	// Get the store as a list
	store := ctx.KVStore(k.storeKey)
	list := list.NewList(k.cdc, store)

	blockTimestamp := ctx.BlockHeader().Time.Unix()

	// Create the record for the block
	blockRecord := &BlockRecord{
		block,
		blockTimestamp,
	}

	formatted, err := json.Marshal(blockRecord)
	if err != nil {
    panic(err)
  }

	index := list.Len()

	// Push in the list with string format
	list.Push(string(formatted))

	// Create the index object to send
	blockIndex := &BlockIndex{
		index,
		blockTimestamp,
	}

	formattedIndex, err := json.Marshal(blockIndex)
	if err != nil {
    panic(err)
  }

	return formattedIndex
}
