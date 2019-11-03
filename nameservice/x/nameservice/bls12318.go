package nameservice

import (

	"github.com/phoreproject/bls/g2pubs"
)

type Bls12381PubKey = [96]byte
type Bls12381Signature = [48]byte

// ValidateSignature verifies a signature against a message and a public key
func ValidateSignature(bls13PubKey Bls12381PubKey, signature Bls12381Signature, message []byte) bool {
	sig, err := g2pubs.DeserializeSignature(signature)
	if err != nil {
		return false
	}

	pubKey, err := g2pubs.DeserializePublicKey(bls13PubKey)
	if err != nil {
		return false
	}

	return g2pubs.Verify(message, pubKey, sig)
}

// ValidatePubKey verifies a public key is on the bls curve?
func ValidatePubKey(bls13PubKey Bls12381PubKey) bool {
	_, err := g2pubs.DeserializePublicKey(bls13PubKey)
	if err != nil {
		return false
	}
	return true
}
