package requestchain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "requestchain" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgAppendBlock:
			return handleAppendBlock(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized requestchain Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleAppendBlock(ctx sdk.Context, keeper Keeper, msg MsgAppendBlock) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}
	keeper.SetName(ctx, msg.Name, msg.Value)

	blockHash := keeper.AppendBlock(ctx, msg.Block)
	if blockHash == "" {
		return sdk.ErrUnauthorized("Block already exists").Result()
	}

	return sdk.Result{}
}
