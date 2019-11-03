package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/sdk-tutorials/nameservice/x/nameservice/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	QueryOrders = "orders"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryOrders:
			return queryOrders(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

// nolint: unparam
func queryOrders(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	// var escrows [][]byte
	var escrows []types.Escrow
	itr := k.GetAllEscrows(ctx)
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		// key := itr.Key()
		escrowBinary := itr.Value()
		// escrows[i] = escrowBinary
		var escrow types.Escrow
		k.cdc.MustUnmarshalBinaryBare(escrowBinary, &escrow)

		escrows = append(escrows, escrow)
	}
	// return k.cdc.MustMarshalBinaryBare(escrows), nil
	res, err := codec.MarshalJSONIndent(k.cdc, escrows)
	if err != nil {
		panic("could not marshal orders query result to JSON")
	}

	return res, nil
}
