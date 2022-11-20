// Copyright 2022 Menabrea Labs Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package language contains types and methods that implement the Marlowe DSL in Go
// See: https://github.com/input-output-hk/marlowe-cardano/tree/main/marlowe/specification
// See: https://github.com/input-output-hk/marlowe-cardano/blob/main/marlowe/src/Language/Marlowe/Core/V1/Semantics/Types.hs
package language_test

import (
	"testing"

	m "github.com/menabrealabs/marlowe/language/v1/core"
)

// Helper function to wrap values in a basic Let contract
func setupLetContract(testVal m.Value) m.Contract {
	return m.Let{
		Name:  "testValue",
		Value: testVal,
		Then:  m.Close,
	}
}

// Helper function to wrap observations in a basic If contract
func setupIfContract(testObs m.Observation) m.Contract {
	return m.If{
		Observe: testObs,
		Then:    m.Close,
		Else:    m.Close,
	}
}

// Tests for stateful value types

func TestTypes_AvailableMoney(t *testing.T) {
	contract := setupLetContract(
		m.AvailableMoney{
			Amount:  m.Ada,
			Account: m.Role{Name: "buyer"},
		},
	)
	assertJson(t, contract,
		`{"let":"testValue","be":{"amount_of_token":{"currency_symbol":"","token_name":""},"in_account":{"role_token":"buyer"}},"then":"close"}`)
}

func TestTypes_ChoiceValue(t *testing.T) {
	contract := setupLetContract(
		m.ChoiceValue{
			Value: m.ChoiceId{
				Name:  "name",
				Owner: m.Role{"buyer"},
			},
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"value_of_choice":{"choice_name":"name","choice_owner":{"role_token":"buyer"}}},"then":"close"}`)
}

func TestTypes_UseValue(t *testing.T) {
	contract := setupLetContract(
		m.UseValue{Value: "value"},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"use_value":"value"},"then":"close"}`)
}

func TestTypes_TimeIntervalStart(t *testing.T) {
	contract := setupLetContract(
		m.TimeIntervalStart,
	)
	assertJson(t, contract, `{"let":"testValue","be":"time_interval_start","then":"close"}`)
}

func TestTypes_TimeIntervalEnd(t *testing.T) {
	contract := setupLetContract(
		m.TimeIntervalEnd,
	)
	assertJson(t, contract, `{"let":"testValue","be":"time_interval_end","then":"close"}`)
}

// Tests for arithmetic value types

func TestTypes_NegValue(t *testing.T) {
	contract := setupLetContract(
		m.NegValue{m.SetConstant("20")},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"negate":20},"then":"close"}`)
}

func TestTypes_AddValue(t *testing.T) {
	contract := setupLetContract(
		m.AddValue{
			Add: m.SetConstant("10"),
			To:  m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"add":10,"and":20},"then":"close"}`)
}

func TestTypes_MulValue(t *testing.T) {
	contract := setupLetContract(
		m.MulValue{
			Multiply: m.SetConstant("10"),
			By:       m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"multiply":10,"times":20},"then":"close"}`)
}

func TestTypes_SubValue(t *testing.T) {
	contract := setupLetContract(
		m.SubValue{
			Subtract: m.SetConstant("10"),
			From:     m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"minus":10,"value":20},"then":"close"}`)
}

func TestTypes_DivValue(t *testing.T) {
	contract := setupLetContract(
		m.DivValue{
			Divide: m.SetConstant("20"),
			By:     m.SetConstant("10"),
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"divide":20,"by":10},"then":"close"}`)
}

// Tests for comparator Observation value types

func TestTypes_ValueGE(t *testing.T) {
	contract := setupIfContract(
		m.ValueGE{
			Value: m.SetConstant("10"),
			Ge:    m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"if":{"value":10,"ge_than":20},"then":"close","else":"close"}`)
}

func TestTypes_ValueGT(t *testing.T) {
	contract := setupIfContract(
		m.ValueGT{
			Value: m.SetConstant("10"),
			Gt:    m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"if":{"value":10,"gt":20},"then":"close","else":"close"}`)
}

func TestTypes_ValueLT(t *testing.T) {
	contract := setupIfContract(
		m.ValueLT{
			Value: m.SetConstant("10"),
			Lt:    m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"if":{"value":10,"lt":20},"then":"close","else":"close"}`)
}

func TestTypes_ValueLE(t *testing.T) {
	contract := setupIfContract(
		m.ValueLE{
			Value: m.SetConstant("10"),
			Le:    m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"if":{"value":10,"le_than":20},"then":"close","else":"close"}`)
}

func TestTypes_ValueEQ(t *testing.T) {
	contract := setupIfContract(
		m.ValueEQ{
			Value: m.SetConstant("10"),
			Eq:    m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"if":{"value":10,"equal_to":20},"then":"close","else":"close"}`)
}

// Tests for connective Observation value types

func TestTypes_NotObs(t *testing.T) {
	contract := setupIfContract(
		m.NotObs{
			m.ValueEQ{
				Value: m.SetConstant("10"),
				Eq:    m.SetConstant("20"),
			},
		},
	)
	assertJson(t, contract, `{"if":{"not":{"value":10,"equal_to":20}},"then":"close","else":"close"}`)
}

func TestTypes_AndObs(t *testing.T) {
	contract := setupIfContract(
		m.AndObs{
			m.ValueEQ{
				Value: m.SetConstant("10"),
				Eq:    m.SetConstant("20"),
			},
			m.ValueEQ{
				Value: m.SetConstant("10"),
				Eq:    m.SetConstant("20"),
			},
		},
	)
	assertJson(t, contract, `{"if":{"both":{"value":10,"equal_to":20},"and":{"value":10,"equal_to":20}},"then":"close","else":"close"}`)
}

func TestTypes_OrObs(t *testing.T) {
	contract := setupIfContract(
		m.OrObs{
			m.ValueEQ{
				Value: m.SetConstant("10"),
				Eq:    m.SetConstant("20"),
			},
			m.ValueEQ{
				Value: m.SetConstant("10"),
				Eq:    m.SetConstant("20"),
			},
		},
	)
	assertJson(t, contract, `{"if":{"either":{"value":10,"equal_to":20},"or":{"value":10,"equal_to":20}},"then":"close","else":"close"}`)
}

func TestTypes_ChoseSomething(t *testing.T) {
	contract := setupIfContract(
		m.ChoseSomething{
			m.ChoiceId{
				Owner: m.Role{Name: "role"},
				Name:  "name",
			},
		},
	)
	assertJson(t, contract, `{"if":{"chose_something_for":{"choice_name":"name","choice_owner":{"role_token":"role"}}},"then":"close","else":"close"}`)
}

func TestTypes_TrueObs(t *testing.T) {
	contract := setupIfContract(m.TrueObs)
	assertJson(t, contract, `{"if":true,"then":"close","else":"close"}`)
}

func TestTypes_FalseObs(t *testing.T) {
	contract := setupIfContract(m.FalseObs)
	assertJson(t, contract, `{"if":false,"then":"close","else":"close"}`)
}
