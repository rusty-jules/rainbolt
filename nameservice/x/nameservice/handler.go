package nameservice

import (
	"fmt"

	"github.com/cosmos/sdk-tutorials/nameservice/x/nameservice/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
		case MsgCreateOrder:
			return handleMsgCreateOrder(ctx, keeper, msg)
		case MsgFillOrder:
			return handleMsgFillOrder(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                      // return
}

// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) sdk.Result {
	// Checks if the the bid price is greater than the price paid by the current owner
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}

// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg MsgDeleteName) sdk.Result {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return types.ErrNameDoesNotExist(types.DefaultCodespace).Result()
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return sdk.ErrUnauthorized("Incorrect Owner").Result()
	}

	keeper.DeleteWhois(ctx, msg.Name)
	return sdk.Result{}
}

// Handle a message to create an order
func handleMsgCreateOrder(ctx sdk.Context, keeper Keeper, msg MsgCreateOrder) sdk.Result {
	// 1. Check if the given BlsPubKey is valid
	// 2. Check that the Merchant is only putting up a single denomination
	// 3. Check if the Merchant already has an escrow account (should get deprecated)
	// 4. Check if the Merchant has the necessary funds to escrow
	// 5. Store Merchant funds in the KV
	// var sigBytes Bls12381PubKey
	// copy(sigBytes[:], msg.Owner)
	// if !ValidatePubKey(sigBytes) {
	// 	return sdk.ErrInvalidPubKey("Invalid BLS public key").Result()
	// }

	if msg.Amount.Len() != 1 {
		return sdk.ErrInternal("Incorrect number of denominations. Must be 1").Result()
	}

	if keeper.IsEscrowPresent(ctx, msg.Merchant.String()) {
		return sdk.ErrInternal("Merchant already has one escrow. Currently only one is supported at a time.").Result()
	}

	coins, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Merchant, msg.Amount)
	if err != nil {
		return sdk.ErrInsufficientCoins("Merchant does not have enough coins to escrow").Result()
	}

	keeper.SetEscrow(ctx, msg.Merchant.String(), Escrow{
		Merchant:     msg.Merchant,
		ChannelState: msg.ChannelState,
		ChannelToken: msg.ChannelToken,
		Amount:       coins,
		Filled:       false,
	})

	return sdk.Result{}
}

// Handle a message to fill an order
func handleMsgFillOrder(ctx sdk.Context, keeper Keeper, msg MsgFillOrder) sdk.Result {
	// 1. Check if the given BlsPubKey is valid
	// 2. Check if the escrow account exists
	// 3. Check if the order has already been filled
	// 4. Check that the amount and denomination put up by the Customer is correct
	// 5. TODO Some stuff to store the correct keys in the escrow for claim
	// 6. Store Customer funds in the KV
	// var sigBytes Bls12381PubKey
	// copy(sigBytes[:], msg.Owner)
	// if !ValidatePubKey(sigBytes) {
	// 	return sdk.ErrInvalidPubKey("Invalid BLS public key").Result()
	// }

	if msg.Amount.Len() != 1 {
		return sdk.ErrInternal("Incorrect number of denominations. Must be 1").Result()
	}

	if !keeper.IsEscrowPresent(ctx, msg.Merchant.String()) {
		return sdk.ErrInternal("Order does not exist. Merchant Escrow not found").Result()
	}

	escrow := keeper.GetEscrow(ctx, msg.Merchant.String())

	if escrow.Filled {
		return sdk.ErrInternal("Order has already been filled").Result()
	}

	if !msg.Amount.IsEqual(escrow.Amount) {
		if !msg.Amount[0].Amount.Equal(escrow.Amount[0].Amount) {
			return sdk.ErrInvalidCoins(fmt.Sprintf("Incorrect amount. Should be %s", escrow.Amount[0].String())).Result()
		}
		return sdk.ErrInvalidCoins(fmt.Sprintf("Incorrect coins. Escrow has %s", escrow.Amount.String())).Result()
	}

	// ---- TODO MuliSig Shizen ----

	coins, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Customer, msg.Amount)
	if err != nil {
		return sdk.ErrInsufficientCoins("Customer does not have enough coins to fill order").Result()
	}

	// This SetEscrow could be done using the already-queried-for-escrow for one less read/deserialize, but whatever
	keeper.SetCustomer(ctx, msg.Merchant.String(), msg.Customer, msg.WalletCommit, coins)
	return sdk.Result{}
}
