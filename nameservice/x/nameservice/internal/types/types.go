package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}

type Bls12381PubKey = [96]byte
type Bls12381Signature = [48]byte

// Escrow is a struct that contains coins held in escrow that will be released
// with a time-delay (for disputes) upon the receipt of a valid Bls12-318 signature
type Escrow struct {
	// This address is needed to verify the Tx that created this escrow, but may not need to be stored by it (unless required to route funds on claim)
	Merchant sdk.AccAddress `json:"merchant"`
	Customer sdk.AccAddress `json:"customer"`
	// A Bls12-318 PublicKey
	ChannelState string    `json:"channelState"`
	ChannelToken string    `json:"channelToken"`
	WalletCommit []byte    `json:"walletState"`
	Amount       sdk.Coins `json:"amount"`
	Filled       bool      `json:"filled"`
	// Denom string `json:"denom"` // stored in sdk.Coin
}

// NewEscrow returns a new Escrow account
func NewEscrow() Escrow {
	return Escrow{
		// Merchant: make([]byte, 96, 96),
		Filled: false,
	}
}

// implement fmt.String
func (e Escrow) String() string {
	return strings.TrimSpace(fmt.Sprintf(
		`Merchant: %s
		Customer: %s
		ChannelState: %s
		ChannelToken: %s
		WalletCommit: %s
		Amount: %s
		Filled: %d`,
		e.Merchant, e.Customer, e.ChannelState, e.ChannelToken, e.WalletCommit, e.Amount, e.Filled,
	))
}
