package requestchain

import (
	"crypto/sha256"
	"github.com/cosmos/cosmos-sdk/codec"
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
func (k Keeper) GetBlock(ctx sdk.Context, hash string) string {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(hash)) {
		return []byte{}
	}

	return store.Get([]byte(hash))
}

// Appends a block of data and returns the hash as an index
func (k Keeper) AppendBlock(ctx sdk.Context, block string) string {
	// Compute sha256 hash of the block
	blockHash := sha256.Sum256([]byte(block))

	store := ctx.KVStore(k.storeKey)
	if !store.Has(blockHash) {
		return ""
	}
	store.Set(blockHash, []byte(block))
	return string(blockHash)
}
