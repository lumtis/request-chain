package requestchain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// "github.com/tendermint/tendermint/crypto"
	// "github.com/tendermint/tendermint/crypto/multisig"
)

// auth.NewAnteHandler(ak, sk, sgc);

// Custom anteHandler to consider fee from block size
func CustomAnteHandler(ak auth.AccountKeeper, supplyKeeper authtypes.SupplyKeeper) sdk.AnteHandler {


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

		// Compute cost
		params := ak.GetParams(newCtx)
		cost := params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes()))

		// Defer
		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case sdk.ErrorOutOfGas:
					log := fmt.Sprintf("out of gas in location: %v", rType.Descriptor)
					res = sdk.ErrOutOfGas(log).Result()
					res.GasWanted = 200000
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					abort = true
				default:
					panic(r)
				}
			}
		}()

    // Consume gas
		newCtx.GasMeter().ConsumeGas(cost, "txSize")

		// Get account
		signerAcc, res := auth.GetSignerAcc(newCtx, ak, stdTx.GetSigners()[0])
		if !res.IsOK() {
			return newCtx, res, true
		}

		// Verify signature
		isGenesis := newCtx.BlockHeight() == 0
		signBytes := auth.GetSignBytes(newCtx.ChainID(), stdTx, signerAcc, isGenesis)
		stdSigs := stdTx.GetSignatures()

		// Set public key to signer if it doesn't exist
		pubKey, res := auth.ProcessPubKey(signerAcc, stdSigs[0], sim)
		if !res.IsOK() {
			return newCtx, res, true
		}
		err := signerAcc.SetPubKey(pubKey)
		if err != nil {
			return newCtx, sdk.ErrInternal("failed to set PubKey on signer account").Result(), true
		}

		if !sim && !pubKey.VerifyBytes(signBytes, stdSigs[0].Signature) {
			return newCtx, sdk.ErrUnauthorized("signature verification failed").Result(), true
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

		// Verify the account has enough funds to pay for fee
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

		// The first signer pays the transaction fees
		res = auth.DeductFees(supplyKeeper, newCtx, signerAcc, fee)
		if !res.IsOK() {
			return newCtx, res, true
		}

		// Send coins to fee collector
		err = supplyKeeper.SendCoinsFromAccountToModule(newCtx, signerAcc.GetAddress(), authtypes.FeeCollectorName, fee)
		if err != nil {
			return newCtx, sdk.ErrInternal(fmt.Sprintf("insufficient spendable funds to pay for fee")).Result(), true
		}

		// Increment sequence
		errSequence := signerAcc.SetSequence(signerAcc.GetSequence() + 1)
		if errSequence != nil {
			return newCtx, sdk.ErrInternal(fmt.Sprintf("impossible to increment sequence")).Result(), true
		}

		ak.SetAccount(newCtx, signerAcc)

		return newCtx, sdk.Result{GasWanted: cost}, false
	}
}
