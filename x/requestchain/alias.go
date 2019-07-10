package requestchain

import (
	"github.com/ltacker/request-chain/x/requestchain/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgAppendBlock = types.NewMsgAppendBlock
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	MsgAppendBlock      = types.MsgAppendBlock
	QueryGetBlock      = types.QueryGetBlock
)
