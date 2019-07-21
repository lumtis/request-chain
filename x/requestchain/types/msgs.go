package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgAppendBlock defines the AppendBlock message
type MsgAppendBlock struct {
	Block  string         `json:"block"`
	Signer sdk.AccAddress `json:"signer"`
}

// NewMsgBuyName is the constructor function for MsgBuyName
func NewMsgAppendBlock(block string, signer sdk.AccAddress) MsgAppendBlock {
	return MsgAppendBlock{
		Block:  block,
		Signer: signer,
	}
}

// Route should return the name of the module
func (msg MsgAppendBlock) Route() string { return RouterKey }

// Type should return the action
func (msg MsgAppendBlock) Type() string { return "append_block" }

// ValidateBasic runs stateless checks on the message
func (msg MsgAppendBlock) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if len(msg.Block) == 0 {
		return sdk.ErrUnknownRequest("Block cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgAppendBlock) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgAppendBlock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
