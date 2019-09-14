package requestchain

import (
	"fmt"
	"strconv"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the requestchain Querier
const (
	QueryGetBlockName = "getblock"
	QueryGetBlockCount = "blockcount"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGetBlockName:
			return queryGetBlock(ctx, path[1:], req, keeper)
		case QueryGetBlockCount:
			return queryGetBlockCount(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown requestchain query endpoint %v", path[0]))
		}
	}
}

// nolint: unparam
func queryGetBlock(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	blockIndex, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
    panic(err)
  }

	res := keeper.GetBlock(ctx, blockIndex)

	if len(res) == 0 {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("Index doesn't exist %v", blockIndex))
	}

	return res, nil
}

// nolint: unparam
func queryGetBlockCount(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	return []byte(fmt.Sprint(keeper.GetBlockCount(ctx))), nil
}
