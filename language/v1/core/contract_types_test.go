package language_test

import (
	"math/big"
	"testing"

	m "github.com/menabrealabs/marlowe/language/v1/core"
)

func TestTypes_CloseContract(t *testing.T) {
	contract := m.Close
	assertJson(t, contract, `"close"`)
}

func TestTypes_LetContract(t *testing.T) {
	// Should generate JSON: {"then":"close","let":"Number","be":1}
	contract := m.Let{
		Name:  "Number",
		Value: m.SetConstant("1"),
		Then:  m.Close,
	}

	assertJson(t, contract, `{"let":"Number","be":1,"then":"close"}`)
}

func TestTypes_IfContract(t *testing.T) {
	// Should generate JSON: {"then":"close","if":{"value":1,"gt":0},"else":"close"}
	contract := m.If{
		Observe: m.ValueGT{
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
		Observe: m.ValueLT{
			Value: m.SetConstant("0"),
			Lt:    m.SetConstant("1"),
		},
		Then: m.Close,
	}

	assertJson(t, contract, `{"assert":{"value":0,"lt":1},"then":"close"}`)
}

func TestTypes_PayContract(t *testing.T) {
	// Should generate JSON:
	// {"token":{"token_name":"","currency_symbol":""},"to":{"party":{"role_token":"creditor"}},"then":"close","pay":5000000,"from_account":{"role_token":"debtor"}}

	contract := m.Pay{
		From:  m.Role{"debtor"},
		To:    m.Payee{m.Role{"creditor"}},
		Token: m.Ada,
		Pay:   m.Constant(*big.NewInt(5_000_000)),
		Then:  m.Close,
	}

	assertJson(t, contract, `{"from_account":{"role_token":"debtor"},"to":{"Party":{"role_token":"creditor"}},"token":{"currency_symbol":"","token_name":""},"pay":5000000,"then":"close"}`)
}

func TestTypes_WhenContract(t *testing.T) {
	// Should generate JSON:
	// {"when":[{"then":"close","case":{"for_choice":{"choice_owner":{"role_token":"creditor"},"choice_name":"option"},"choose_between":[{"to":2,"from":1}]}}],"timeout_continuation":"close","timeout":1668250824063}

	contract := m.When{
		Cases: []m.Case{
			{
				Action: m.Choice{
					ChoiceId: m.ChoiceId{
						Name:  "option",
						Owner: m.Role{"creditor"},
					},
					Bounds: []m.Bound{
						{
							Upper: 3,
							Lower: 2,
						},
					},
				},
				Then: m.Close,
			},
		},
		Timeout: m.POSIXTime(1666078977926),
		Then:    m.Close,
	}

	assertJson(t, contract, `{"when":[{"case":{"for_choice":{"choice_name":"option","choice_owner":{"role_token":"creditor"}},"choose_between":[{"from":3,"to":2}]},"then":"close"}],"timeout":1666078977926,"timeout_continuation":"close"}`)
}
