package language

import (
	"github.com/btcsuite/btcutil/bech32"
)

// Validate that an address will decode from a Bech32 encoding to a valid ledger address.
func (a Address) ValidateEncoding() error {
	_, _, err := bech32.Decode(string(a))

	if err != nil {
		return err
	}

	return nil
}
