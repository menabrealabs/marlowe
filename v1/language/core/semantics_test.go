package language_test

import (
	"testing"

	lang "github.com/menabrealabs/marlowe/v1/language/core"
)

func TestAddress_ValidateEncoding_ShouldPass(t *testing.T) {
	// Test vectors specified in BIP-173:
	// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki#Test_vectors
	testVectors := []string{
		"A12UEL5L",
		"a12uel5l",
		"an83characterlonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1tt5tgs",
		"abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
		"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j",
		"split1checkupstagehandshakeupstreamerranterredcaperred2y9e3w",
		"?1ezyfcl",
	}

	for _, bech32 := range testVectors {
		addr := lang.Address(bech32)
		err := addr.ValidateEncoding()

		if err != nil {
			t.Error(err)
		}
	}
}

func TestAddress_ValidateEncoding_ShouldFail(t *testing.T) {
	// Test vectors specified in BIP-173:
	// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki#Test_vectors
	testVectors := []string{
		"an84characterslonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1569pvx", // overall max length exceeded
		"pzry9x0s0muk",  // No separator character
		"1pzry9x0s0muk", // Empty HRP
		"x1b4n0q5v",     // Invalid data character
		"li1dgmt3",      // Too short checksum
		"A1G7SGD8",      // Checksum calculated with uppercase form of HRP
		"10a06t8",       // Empty HRP
		"1qzzfhee",      // Empty HRP
	}

	for _, bech32 := range testVectors {
		addr := lang.Address(bech32)
		err := addr.ValidateEncoding()

		t.Log("Expected error: ", err)

		if err == nil {
			t.Error("Invalid address '", bech32, "' should have failed validation but didn't.")
		}
	}
}
