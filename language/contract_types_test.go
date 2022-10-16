package language_test

import (
	"encoding/json"
	"math/big"
	"testing"

	m "github.com/menabrealabs/marlowe/language"
)

func TestTypes_CloseContract(t *testing.T) {
	contract := m.Close
	assertJson(t, contract, `"close"`)
}

func TestTypes_LetContract(t *testing.T) {
	// Should generate JSON: {"then":"close","let":"Number","be":1}
	contract := m.Let{
		"Number",
		m.Constant(*big.NewInt(1)),
		m.Close,
	}

	assertJson(t, contract, `{"let":"Number","be":1,"then":"close"}`)
}

func TestTypes_IfContract(t *testing.T) {
	// Should generate JSON: {"then":"close","if":{"value":1,"gt":0},"else":"close"}
	contract := m.If{
		m.ValueGT{m.Constant(*big.NewInt(1)), m.Constant(*big.NewInt(0))},
		m.Close,
		m.Close,
	}

	assertJson(t, contract, `{"if":{"value":1,"gt":0},"then":"close","else":"close"}`)
}

func TestTypes_AssertContract(t *testing.T) {
	// Should generate JSON: {"then":"close","assert":{"value":0,"lt":1}}
	contract := m.Assert{
		m.ValueLT{m.Constant(*big.NewInt(0)), m.Constant(*big.NewInt(1))},
		m.Close,
	}

	assertJson(t, contract, `{"assert":{"value":0,"lt":1},"then":"close"}`)
}

func TestTypes_PayContract(t *testing.T) {
	// Should generate JSON:
	// {"token":{"token_name":"","currency_symbol":""},"to":{"party":{"role_token":"creditor"}},"then":"close","pay":5000000,"from_account":{"role_token":"debtor"}}

	contract := m.Pay{
		m.Role{"debitor"},
		m.Payee{m.Role{"creditor"}},
		m.Ada,
		m.Constant(*big.NewInt(50000000)),
		m.Close,
	}

	assertJson(t, contract, `{"from_account":{"role_token":"debitor"},"to":{"Party":{"role_token":"creditor"}},"token":{"currency_symbol":"","token_name":""},"pay":5000000,"then":"close"}`)
}

func assertJson(t *testing.T, contract m.Contract, target string) {
	jbytes, err := json.Marshal(contract)
	if err != nil {
		t.Error(err)
	}

	if string(jbytes) != target {
		t.Error("Incorrect JSON format: ", string(jbytes))
	}

	t.Logf("Marshalled JSON: %v", string(jbytes))
}
