package requestchain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	// tmcrypto "github.com/tendermint/tendermint/crypto"
)

// Custom anteHandler to consider fee from block size
func CustomAnteHandler(ak auth.AccountKeeper, fck auth.FeeCollectionKeeper) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {

		// switch castTx := tx.(type) {
		// 	case auth.StdTx:
		// 		return sdkAnteHandler(ctx, ak, fck, castTx, sim)
    //
		// 	default:
		// 		return ctx, sdk.ErrInternal(fmt.Sprintf("transaction type invalid: %T", tx)).Result(), true
		// }

		return ctx, sdk.ErrInternal(fmt.Sprintf("transaction type invalid: %T", tx)).Result(), true
	}
}

// func sdkAnteHandler(
// 	ctx sdk.Context, ak auth.AccountKeeper, fck auth.FeeCollectionKeeper, stdTx auth.StdTx, sim bool,
// ) (newCtx sdk.Context, res sdk.Result, abort bool) {
//
// 	// Ensure that the provided fees meet a minimum threshold for the validator,
// 	// if this is a CheckTx. This is only for local mempool purposes, and thus
// 	// is only ran on check tx.
// 	if ctx.IsCheckTx() && !sim {
// 		res := auth.EnsureSufficientMempoolFees(ctx, stdTx)
// 		if !res.IsOK() {
// 			return newCtx, res, true
// 		}
// 	}
//
// 	newCtx = auth.SetGasMeter(sim, ctx, stdTx)
//
// 	// AnteHandlers must have their own defer/recover in order for the BaseApp
// 	// to know how much gas was used! This is because the GasMeter is created in
// 	// the AnteHandler, but if it panics the context won't be set properly in
// 	// runTx's recover call.
// 	defer func() {
// 		if r := recover(); r != nil {
// 			switch rType := r.(type) {
// 			case sdk.ErrorOutOfGas:
// 				log := fmt.Sprintf("out of gas in location: %v", rType.Descriptor)
// 				res = sdk.ErrOutOfGas(log).Result()
// 				res.GasWanted = stdTx.Fee.Gas
// 				res.GasUsed = newCtx.GasMeter().GasConsumed()
// 				abort = true
// 			default:
// 				panic(r)
// 			}
// 		}
// 	}()
//
// 	if err := stdTx.ValidateBasic(); err != nil {
// 		return newCtx, err.Result(), true
// 	}
//
// 	newCtx.GasMeter().ConsumeGas(memoCostPerByte*sdk.Gas(len(stdTx.GetMemo())), "memo")
//
// 	signerAccs, res := auth.GetSignerAccs(newCtx, ak, stdTx.GetSigners())
// 	if !res.IsOK() {
// 		return newCtx, res, true
// 	}
//
// 	// the first signer pays the transaction fees
// 	if !stdTx.Fee.Amount.IsZero() {
// 		signerAccs[0], res = auth.DeductFees(signerAccs[0], stdTx.Fee)
// 		if !res.IsOK() {
// 			return newCtx, res, true
// 		}
//
// 		fck.AddCollectedFees(newCtx, stdTx.Fee.Amount)
// 	}
//
// 	isGenesis := ctx.BlockHeight() == 0
// 	signBytesList := auth.GetSignBytesList(newCtx.ChainID(), stdTx, signerAccs, isGenesis)
// 	stdSigs := stdTx.GetSignatures()
//
// 	for i := 0; i < len(stdSigs); i++ {
// 		// check signature, return account with incremented nonce
// 		signerAccs[i], res = processSig(newCtx, signerAccs[i], stdSigs[i], signBytesList[i], sim)
// 		if !res.IsOK() {
// 			return newCtx, res, true
// 		}
//
// 		ak.SetAccount(newCtx, signerAccs[i])
// 	}
//
// 	return newCtx, sdk.Result{GasWanted: stdTx.Fee.Gas}, false
// }
//
// // processSig verifies the signature and increments the nonce. If the account
// // doesn't have a pubkey, set it.
// func processSig(
// 	ctx sdk.Context, acc auth.Account, sig auth.StdSignature, signBytes []byte, sim bool,
// ) (updatedAcc auth.Account, res sdk.Result) {
//
// 	pubKey, res := auth.ProcessPubKey(acc, sig, sim)
// 	if !res.IsOK() {
// 		return nil, res
// 	}
//
// 	err := acc.SetPubKey(pubKey)
// 	if err != nil {
// 		return nil, sdk.ErrInternal("failed to set PubKey on signer account").Result()
// 	}
//
// 	consumeSigGas(ctx.GasMeter(), pubKey)
// 	if !sim && !pubKey.VerifyBytes(signBytes, sig.Signature) {
// 		return nil, sdk.ErrUnauthorized("signature verification failed").Result()
// 	}
//
// 	err = acc.SetSequence(acc.GetSequence() + 1)
// 	if err != nil {
// 		return nil, sdk.ErrInternal("failed to set account nonce").Result()
// 	}
//
// 	return acc, res
// }
//
// func consumeSigGas(meter sdk.GasMeter, pubkey tmcrypto.PubKey) {
// 	switch pubkey.(type) {
// 	case crypto.PubKeySecp256k1:
// 		meter.ConsumeGas(secp256k1VerifyCost, "ante verify: secp256k1")
// 	default:
// 		panic("Unrecognized signature type")
// 	}
// }
