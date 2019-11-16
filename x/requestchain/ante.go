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
func CustomAnteHandler(ak auth.AccountKeeper) sdk.AnteHandler {


	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {
		// TODO: implements SigVerifiableTx
		// sigTx, ok := tx.(SigVerifiableTx)
		// if !ok {
		// 	fmt.Println("invalid tx type") // sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
		// 	return ctx, sdk.Result{}, true
		// }


		params := ak.GetParams(ctx)
		fmt.Println("tx size to consume: ", sdk.Gas(len(ctx.TxBytes())))
		ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(ctx.TxBytes())), "txSize")

		// TODO: Implement simulate
		// // simulate gas cost for signatures in simulate mode
		// if sim {
		// 	// in simulate mode, each element should be a nil signature
		// 	sigs := sigTx.GetSignatures()
		// 	for i, signer := range sigTx.GetSigners() {
		// 		// if signature is already filled in, no need to simulate gas cost
		// 		if sigs[i] != nil {
		// 			continue
		// 		}
		// 		acc := ak.GetAccount(ctx, signer)

		// 		var pubkey crypto.PubKey
		// 		// use placeholder simSecp256k1Pubkey if sig is nil
		// 		if acc == nil || acc.GetPubKey() == nil {
		// 			pubkey = simSecp256k1Pubkey
		// 		} else {
		// 			pubkey = acc.GetPubKey()
		// 		}
		// 		// use stdsignature to mock the size of a full signature
		// 		simSig := types.StdSignature{
		// 			Signature: simSecp256k1Sig[:],
		// 			PubKey:    pubkey,
		// 		}
		// 		sigBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(simSig)
		// 		cost := sdk.Gas(len(sigBz) + 6)

		// 		// If the pubkey is a multi-signature pubkey, then we estimate for the maximum
		// 		// number of signers.
		// 		if _, ok := pubkey.(multisig.PubKeyMultisigThreshold); ok {
		// 			cost *= params.TxSigLimit
		// 		}

		// 		ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*cost, "txSize")
		// 	}
		// }

		return ctx, sdk.Result{}, false
	}
}
