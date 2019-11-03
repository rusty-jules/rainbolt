package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/sdk-tutorials/nameservice/x/nameservice/internal/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	CoinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetWhois(ctx sdk.Context, name string) types.Whois {
	store := ctx.KVStore(k.storeKey)
	if !k.IsNamePresent(ctx, name) {
		return types.NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois types.Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois types.Whois) {
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

// Deletes the entire Whois metadata struct for a name
func (k Keeper) DeleteWhois(ctx sdk.Context, name string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(name))
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

// GetPrice - gets the current price of a name
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

// Check if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// SetEscrow an escrow account. Currently uses the sender address (which means only one escrow account per user) as a key. This should be changed (and obfuscated) in a clever way
func (k Keeper) SetEscrow(ctx sdk.Context, senderAddress string, escrow types.Escrow) {
	if escrow.Amount.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(senderAddress), k.cdc.MustMarshalBinaryBare(escrow))
}

// Get the Escrow data for a sender address
// func (k Keeper) GetEscrow(ctx sdk.Context, senderAddress string) Escrow {
// 	store := ctx.KVStore(k.storeKey)
// 	if !k.IsEscrowPresent(ctx, senderAddress) {
// 		return NewEscrow() // FIXME, maybe throw error?
// 	}
// }

// IsEscrowPresent checks if an Escrow account exists for the given sender address
func (k Keeper) IsEscrowPresent(ctx sdk.Context, senderAddress string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(senderAddress))
}

// IsEscrowFilled checks if an Escrow account has been filled
func (k Keeper) IsEscrowFilled(ctx sdk.Context, senderAddress string) bool {
	return k.GetEscrow(ctx, senderAddress).Filled
}

// GetEscrow returns an escrow account given a key
func (k Keeper) GetEscrow(ctx sdk.Context, senderAddress string) types.Escrow {
	store := ctx.KVStore(k.storeKey)
	if !k.IsEscrowPresent(ctx, senderAddress) {
		return types.NewEscrow()
	}
	binary := store.Get([]byte(senderAddress))
	var escrow types.Escrow
	k.cdc.MustUnmarshalBinaryBare(binary, &escrow)
	return escrow
}

// GetEscrowSize returns an escrows contract size
func (k Keeper) GetEscrowSize(ctx sdk.Context, senderAddress string) sdk.Int {
	escrow := k.GetEscrow(ctx, senderAddress)
	if escrow.Filled {
		return escrow.Amount[0].Amount.QuoRaw(2)
	}
	return escrow.Amount[0].Amount
}

// GetEscrowDenom returns an escrows denomination
func (k Keeper) GetEscrowDenom(ctx sdk.Context, senderAddress string) string {
	return k.GetEscrow(ctx, senderAddress).Amount[0].Denom
}

// SetCustomer adds a customer to the escrow account
func (k Keeper) SetCustomer(ctx sdk.Context, senderAddress string, customer sdk.AccAddress, walletCommit []byte, coins sdk.Coins) {
	escrow := k.GetEscrow(ctx, senderAddress)
	escrow.Customer = customer
	escrow.WalletCommit = walletCommit
	escrow.Amount.Add(coins)
	escrow.Filled = true
	k.SetEscrow(ctx, senderAddress, escrow)
}

// GetAllEscrows lets you see all "orders" on chain
func (k Keeper) GetAllEscrows(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
