package language_test

import (
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
		Name:     "Number",
		Value:    m.SetConstant("1"),
		Continue: m.Close,
	}

	assertJson(t, contract, `{"let":"Number","be":1,"then":"close"}`)
}

func TestTypes_IfContract(t *testing.T) {
	// Should generate JSON: {"then":"close","if":{"value":1,"gt":0},"else":"close"}
	contract := m.If{
		Observation: m.ValueGT{
			Value: m.SetConstant("1"),
			Gt:    m.SetConstant("0"),
		},
		Then: m.Close,
		Else: m.Close,
	}

	assertJson(t, contract, `{"if":{"value":1,"gt":0},"then":"close","else":"close"}`)
}

func TestTypes_AssertContract(t *testing.T) {
	// Should generate JSON: {"then":"close","assert":{"value":0,"lt":1}}
	contract := m.Assert{
		Observation: m.ValueLT{
			Value: m.SetConstant("0"),
			Lt:    m.SetConstant("1"),
		},
		Continue: m.Close,
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
		m.Constant(*big.NewInt(5000000)),
		m.Close,
	}

	assertJson(t, contract, `{"from_account":{"role_token":"debitor"},"to":{"Party":{"role_token":"creditor"}},"token":{"currency_symbol":"","token_name":""},"pay":5000000,"then":"close"}`)
}

func TestTypes_WhenContract(t *testing.T) {
	// Should generate JSON:
	// {"when":[{"then":"close","case":{"for_choice":{"choice_owner":{"pk_hash":"0000000000000000000000000000000000000000000000000000000000000000"},"choice_name":"option"},"choose_between":[{"to":2,"from":1}]}}],"timeout_continuation":"close","timeout":1666078977926}

	contract := m.When{
		[]m.Case{},
		1666078977926,
		m.Close,
	}

	assertJson(t, contract, `{"from_account":{"role_token":"debitor"},"to":{"Party":{"role_token":"creditor"}},"token":{"currency_symbol":"","token_name":""},"pay":5000000,"then":"close"}`)
}
