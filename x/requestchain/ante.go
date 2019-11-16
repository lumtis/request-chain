package requestchain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// Custom anteHandler to consider fee from block size
func CustomAnteHandler(	ak auth.AccountKeeper,
												fck auth.FeeCollectionKeeper,
												sk types.SupplyKeeper,
												sgc SignatureVerificationGasConsumer) sdk.AnteHandler {

	// auth.NewAnteHandler(ak, sk, sgc);
	cgtsd := ante.NewConsumeGasForTxSizeDecorator(app.AccountKeeper)
	antehandler := sdk.ChainAnteDecorators(cgtsd)

	return antehandler;
}

func AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}
	params := cgts.ak.GetParams(ctx)
	ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(ctx.TxBytes())), "txSize")

	// simulate gas cost for signatures in simulate mode
	if simulate {
		// in simulate mode, each element should be a nil signature
		sigs := sigTx.GetSignatures()
		for i, signer := range sigTx.GetSigners() {
			// if signature is already filled in, no need to simulate gas cost
			if sigs[i] != nil {
				continue
			}
			acc := cgts.ak.GetAccount(ctx, signer)

			var pubkey crypto.PubKey
			// use placeholder simSecp256k1Pubkey if sig is nil
			if acc == nil || acc.GetPubKey() == nil {
				pubkey = simSecp256k1Pubkey
			} else {
				pubkey = acc.GetPubKey()
			}
			// use stdsignature to mock the size of a full signature
			simSig := types.StdSignature{
				Signature: simSecp256k1Sig[:],
				PubKey:    pubkey,
			}
			sigBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(simSig)
			cost := sdk.Gas(len(sigBz) + 6)

			// If the pubkey is a multi-signature pubkey, then we estimate for the maximum
			// number of signers.
			if _, ok := pubkey.(multisig.PubKeyMultisigThreshold); ok {
				cost *= params.TxSigLimit
			}

			ctx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*cost, "txSize")
		}
	}

	return next(ctx, tx, simulate)
}
