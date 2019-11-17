package requestchain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	// "github.com/cosmos/cosmos-sdk/x/auth/types"
	// "github.com/tendermint/tendermint/crypto"
	// "github.com/tendermint/tendermint/crypto/multisig"
)

// auth.NewAnteHandler(ak, sk, sgc);

// Custom anteHandler to consider fee from block size
func CustomAnteHandler(ak auth.AccountKeeper, fck auth.FeeCollectionKeeper) sdk.AnteHandler {


	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {
		// TODO: implements SigVerifiableTx
		// sigTx, ok := tx.(SigVerifiableTx)
		// if !ok {
		// 	fmt.Println("invalid tx type") // sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
		// 	return ctx, sdk.Result{}, true
		// }

		// TODO: Implement simulate
		// https://github.com/cosmos/cosmos-sdk/blob/28347bf5f7369a8ee1673c19be51f723b686b650/x/auth/ante/basic.go

		stdTx, ok := tx.(auth.StdTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be StdTx").Result(), true
		}

		newCtx = auth.SetGasMeter(sim, ctx, 200000)

    // Consume gas
		params := ak.GetParams(newCtx)
		cost := params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes()))
		newCtx.GasMeter().ConsumeGas(cost, "txSize")

		// Get account
		signerAcc, res := auth.GetSignerAcc(newCtx, ak, stdTx.GetSigners()[0])
		if !res.IsOK() {
			return newCtx, res, true
		}

		// Deduct fee
		blockTime := newCtx.BlockHeader().Time
		coins := signerAcc.GetCoins()

		costCoin := sdk.NewInt64Coin("stake", int64(cost))
		fee := sdk.Coins([]sdk.Coin{costCoin})

		if !fee.IsValid() {
      // sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fee)
			return newCtx, sdk.ErrInternal(fmt.Sprintf("invalid fee amount: %s", fee)).Result(), true
		}

		// verify the account has enough funds to pay for fee
		_, hasNeg := coins.SafeSub(fee)
		if hasNeg {
      // sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient funds to pay for fee; %s < %s", coins, fee)
			return newCtx, sdk.ErrInternal(fmt.Sprintf("insufficient funds to pay for fee; %s < %s", coins, fee)).Result(), true
		}

		// Validate the account has enough "spendable" coins as this will cover cases
		// such as vesting accounts.
		spendableCoins := signerAcc.SpendableCoins(blockTime)
		if _, hasNeg := spendableCoins.SafeSub(fee); hasNeg {
			// sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient funds to pay for fee; %s < %s", spendableCoins, fee)
			return newCtx, sdk.ErrInternal(fmt.Sprintf("insufficient spendable funds to pay for fee; %s < %s", spendableCoins, fee)).Result(), true
		}

		// the first signer pays the transaction fees
		signerAcc, res = auth.DeductFees(blockTime, signerAcc, auth.NewStdFee(200000, fee))
		if !res.IsOK() {
			return newCtx, res, true
		}
		fck.AddCollectedFees(newCtx, fee)

		return newCtx, sdk.Result{GasWanted: cost}, false
	}
}
