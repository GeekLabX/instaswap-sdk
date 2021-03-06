package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/jpthor/cosmos-swap/x/swapservice/types"

	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

const (
	restPoolData = "pooldata"
	stakeData    = "stakedata"
	accData      = "accdata"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storePoolData string) {
	r.HandleFunc(fmt.Sprintf("/%s/pools", storePoolData), poolHandler(cdc, cliCtx, storePoolData)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/account/{%s}", storePoolData, accData), accHandler(cdc, cliCtx, storePoolData)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/stake/{%s}", storePoolData, stakeData), stakeHandler(cdc, cliCtx, storePoolData)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/stake", storePoolData), setStakeDataHandler(cdc, cliCtx)).Methods("PUT")
}

type setStakeData struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Name    string       `json:"name"`
	Ticker  string       `json:"ticker"`
	Atom    string       `json:"atom_amount"`
	Token   string       `json:"token_amount"`
}

func setStakeDataHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		var req setStakeData
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgSetStakeData(req.Name, req.Ticker, req.Atom, req.Token, cliCtx.GetFromAddress())
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
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/poolstruct/%s", storePoolData, paramType), nil)
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

func accHandler(cdc *codec.Codec, cliCtx context.CLIContext, storePoolData string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[accData]
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/accstruct/%s", storePoolData, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func stakeHandler(cdc *codec.Codec, cliCtx context.CLIContext, storePoolData string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[stakeData]
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/stakestruct/%s", storePoolData, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}
