package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/sdk-tutorials/nameservice/x/nameservice/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

type createOrderReq struct {
	BaseReq      rest.BaseReq `json:"base_req"`
	Merchant     string       `json:"merchant"`
	ChannelState string       `json:"channelState"`
	ChannelToken string       `json:"channelToken"`
	Amount       string       `json:"amount"`
}

func createOrderHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createOrderReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		addr, err := sdk.AccAddressFromBech32(req.Merchant)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// var owner, commitPubKey [96]byte
		// blsPubKeyBytes := []byte(req.Owner)
		// commitPubKeyBytes := []byte(req.CommitPubKey)
		// copy(owner[:], blsPubKeyBytes)
		// copy(commitPubKey[:], commitPubKeyBytes)

		coins, err := sdk.ParseCoins(req.Amount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Message time
		msg := types.NewMsgCreateOrder(addr, req.ChannelState, req.ChannelToken, coins)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
