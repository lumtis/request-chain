package requestchain

import (
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
	var block string

	if list.Len() <= index {
		return []byte{}
	}

	err := list.Get(index, &block)
	if err != nil {
    panic(err)
  }

	// blockString := fmt.Sprintf("%v", block)

	return []byte(block)
}

// Appends a block of data and returns the hash as an index
func (k Keeper) AppendBlock(ctx sdk.Context, block string) {
	// Get the store as a lsit
	store := ctx.KVStore(k.storeKey)
	list := list.NewList(k.cdc, store)

	list.Push(block) // []byte(block)
	}
