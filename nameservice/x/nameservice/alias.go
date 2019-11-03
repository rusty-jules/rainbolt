package nameservice

import (
	"github.com/cosmos/sdk-tutorials/nameservice/x/nameservice/internal/keeper"
	"github.com/cosmos/sdk-tutorials/nameservice/x/nameservice/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper        = keeper.NewKeeper
	NewQuerier       = keeper.NewQuerier
	NewMsgBuyName    = types.NewMsgBuyName
	NewMsgSetName    = types.NewMsgSetName
	NewMsgDeleteName = types.NewMsgDeleteName

	NewMsgCreateOrder = types.NewMsgCreateOrder
	NewMsgFillOrder   = types.NewMsgFillOrder
	NewEscrow         = types.NewEscrow

	NewWhois      = types.NewWhois
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper          = keeper.Keeper
	MsgSetName      = types.MsgSetName
	MsgBuyName      = types.MsgBuyName
	MsgDeleteName   = types.MsgDeleteName
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames

	MsgCreateOrder = types.MsgCreateOrder
	MsgFillOrder   = types.MsgFillOrder
	QueryResOrder  = types.QueryResOrders
	Escrow         = types.Escrow

	Whois = types.Whois
)
