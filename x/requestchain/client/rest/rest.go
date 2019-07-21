package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ltacker/request-chain/x/requestchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/gorilla/mux"
)

const (
	restName = "rc"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeBlock string) {
	r.HandleFunc(fmt.Sprintf("/%s/blocks", storeBlock), appendBlockHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/getblock/{%s}", storeBlock, restName), getBlockHandler(cliCtx, storeBlock)).Methods("GET")
}

// --------------------------------------------------------------------------------------
// Tx Handler

type appendBlockReq struct {
	BaseReq rest.BaseReq  `json:"base_req"`
	Block   string        `json:"block"`
	Signer sdk.AccAddress `json:"signer"`
}

func appendBlockHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req appendBlockReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgAppendBlock(req.Block, req.Signer)
		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

//--------------------------------------------------------------------------------------
// Query Handlers

func getBlockHandler(cliCtx context.CLIContext, storeBlock string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/get/%s", storeBlock, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
