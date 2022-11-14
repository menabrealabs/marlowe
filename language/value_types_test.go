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

	m "github.com/menabrealabs/marlowe/language"
)

func setupContract(testVal m.Value) m.Contract {
	return m.Let{
		Name:  "testValue",
		Value: testVal,
		Then:  m.Close,
	}
}

func TestTypes_NegValue(t *testing.T) {
	contract := setupContract(
		m.NegValue{m.SetConstant("20")},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"negate":20},"then":"close"}`)
}

func TestTypes_AddValue(t *testing.T) {
	contract := setupContract(
		m.AddValue{
			Add: m.SetConstant("10"),
			To:  m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"add":10,"and":20},"then":"close"}`)
}

func TestTypes_MulValue(t *testing.T) {
	contract := setupContract(
		m.MulValue{
			Multiply: m.SetConstant("10"),
			By:       m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"multiply":10,"times":20},"then":"close"}`)
}

func TestTypes_SubValue(t *testing.T) {
	contract := setupContract(
		m.SubValue{
			Subtract: m.SetConstant("10"),
			From:     m.SetConstant("20"),
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"minus":10,"value":20},"then":"close"}`)
}

func TestTypes_DivValue(t *testing.T) {
	contract := setupContract(
		m.DivValue{
			Divide: m.SetConstant("20"),
			By:     m.SetConstant("10"),
		},
	)
	assertJson(t, contract, `{"let":"testValue","be":{"divide":20,"by":10},"then":"close"}`)
}
