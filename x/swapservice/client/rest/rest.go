package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/jpthor/test/x/swapservice/types"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

const (
	restPoolData = "pooldata"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storePoolData string) {
	r.HandleFunc(fmt.Sprintf("/%s/pools", storePoolData), poolHandler(cdc, cliCtx, storePoolData)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pools", storePoolData), setPoolDataHandler(cdc, cliCtx)).Methods("PUT")
	//r.HandleFunc(fmt.Sprintf("/%s/pools/{%s}", storePoolData, restPoolData), resolvePoolDataHandler(cdc, cliCtx, storePoolData)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/pools/{%s}/poolstruct", storePoolData, restPoolData), whoIsHandler(cdc, cliCtx, storePoolData)).Methods("GET")
}

type setPoolDataReq struct {
	BaseReq      rest.BaseReq `json:"base_req"`
	TokenName    string       `json:"token_name"`
	Ticker       string       `json:"ticker"`
	BalanceAtom  string       `json:"balance_atom"`
	BalanceToken string       `json:"balance_token"`
}

func setPoolDataHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		var req setPoolDataReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgSetPoolData(req.TokenName, req.Ticker, cliCtx.GetFromAddress())
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		clientrest.WriteGenerateStdTxResponse(w, cdc, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func getPoolDataHandler(cdc *codec.Codec, cliCtx context.CLIContext, storePoolData string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restPoolData]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getpool/%s", storePoolData, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func whoIsHandler(cdc *codec.Codec, cliCtx context.CLIContext, storePoolData string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restPoolData]
		log.Printf("Got here")

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/poolstruct/%s", storePoolData, paramType), nil)
		log.Printf("RES: %+v", res)
		log.Printf("Err: %+v", err)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func poolHandler(cdc *codec.Codec, cliCtx context.CLIContext, storePoolData string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pooldatas", storePoolData), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}