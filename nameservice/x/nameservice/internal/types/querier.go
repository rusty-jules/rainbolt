package types

import "strings"

// QueryResResolve Queries Result Payload for a resolve query
type QueryResResolve struct {
	Value string `json:"value"`
}

// implement fmt.Stringer
func (r QueryResResolve) String() string {
	return r.Value
}

// QueryResNames Queries Result Payload for a names query
type QueryResNames []string

// implement fmt.Stringer
func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}

// Query Result Payload for Orders
type QueryResOrders []Escrow

// implement fmt.Stringer
// TODO might need to make this a JSON object?
func (n QueryResOrders) String() string {
	var escrowStrings []string
	for i := 0; i < len(n); i++ {
		escrowStrings[i] = n[i].String()
	}

	return strings.Join(escrowStrings[:], "\n")
}
