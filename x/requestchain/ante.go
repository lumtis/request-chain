package requestchain

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

type ChainAnte struct{}

func (ChainAnte) AnteHandle(ctx Context, tx Tx, simulate bool, next AnteHandler) (newCtx Context, err error {
	fmt.Println(tx)
}


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
