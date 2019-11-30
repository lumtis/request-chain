package rest

import (
	"fmt"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ltacker/request-chain/x/requestchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/gorilla/mux"
)

const (
	restName = "rc"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeBlock string) {
	r.HandleFunc(fmt.Sprintf("/%s/blocks", storeBlock), appendBlockHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/getblock/{%s}", storeBlock, restName), getBlockHandler(cliCtx, storeBlock)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/blockcount", storeBlock), getBlockCountHandler(cliCtx, storeBlock)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/broadcast", storeBlock), broadcastHandler(cliCtx)).Methods("POST")
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

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getblock/%s", storeBlock, paramType), nil)
		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getBlockCountHandler(cliCtx context.CLIContext, storeBlock string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/blockcount", storeBlock), nil)
		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//--------------------------------------------------------------------------------------
// Signature and broadcasting

type MessageBody struct {
	Tx               auth.StdTx `json:"tx"`
	LocalAccountName string     `json:"name"`
	Address					 string			`json:"address"`
	Password         string     `json:"password"`
	ChainID          string     `json:"chain_id"`
	AccountNumber    int64      `json:"account_number"`
}


func broadcastHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m MessageBody

		fmt.Println("Broadcast called")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = cliCtx.Codec.UnmarshalJSON(body, &m)

		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		accountRetriever := authtypes.NewAccountRetriever(cliCtx)
		acc, err := accountRetriever.GetAccount(sdk.AccAddress(m.Address))
		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		seq := acc.GetSequence()

		fmt.Println("Signing transaction")

		// Sign
		stdSignature, err := MakeSignature(m.LocalAccountName, m.Password, auth.StdSignMsg{
			ChainID:       m.ChainID,
			AccountNumber: uint64(m.AccountNumber),
			Sequence:      seq,
			Fee:           m.Tx.Fee,
			Msgs:          m.Tx.GetMsgs(),
			Memo:          m.Tx.GetMemo(),
		})
		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		sigs := []auth.StdSignature{stdSignature}

		signedStdTx := auth.NewStdTx(m.Tx.GetMsgs(), m.Tx.Fee, sigs, m.Tx.GetMemo())

		fmt.Println("Encoding block")

		encoder := utils.GetTxEncoder(cliCtx.Codec)
		txToBroadcast, err := encoder(signedStdTx)
		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("Broadcasting block")

		// broadcast to a Tendermint node
		res, err := cliCtx.BroadcastTxCommit(txToBroadcast)
		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Decode the block record
		blockRecord, err := hex.DecodeString(res.Data)
		if err != nil {
			fmt.Println(err)
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("Block has been broadcasted")

		// Send block record
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(blockRecord)
	}
}

// MakeSignature builds a StdSignature given key name, passphrase, and a StdSignMsg.
func MakeSignature(name, passphrase string, msg auth.StdSignMsg) (sig auth.StdSignature, err error) {
	keybase, err := keys.NewKeyBaseFromHomeFlag()
	if err != nil {
		return
	}
	sigBytes, pubkey, err := keybase.Sign(name, passphrase, msg.Bytes())
	if err != nil {
		return
	}
	return auth.StdSignature{
		PubKey:        pubkey,
		Signature:     sigBytes,
	}, nil
}
