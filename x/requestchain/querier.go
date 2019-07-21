package requestchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the requestchain Querier
const (
	QueryGetBlockName = "getblock"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGetBlockName:
			return queryGetBlock(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown requestchain query endpoint")
		}
	}
}

// nolint: unparam
func queryGetBlock(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	res := keeper.GetBlock(ctx, path[0])

	if len(res) == 0 {
		return []byte{}, sdk.ErrUnknownRequest("Hash doesn't exist")
	}

	return res, nil
}
