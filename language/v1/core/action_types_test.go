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

// Note: the Choice action is already tested in TestTypes_When() in content_types_test.go

func setupWhenContract(a m.Action) m.Contract {
	return m.When{
		Cases: []m.Case{
			{
				Action: a,
				Then:   m.Close,
			},
		},
		Timeout: 1666078977926,
		Then:    m.Close,
	}
}

func TestTypes_Deposit(t *testing.T) {
	contract := setupWhenContract(
		m.Deposit{
			IntoAccount: m.Role{"seller"},
			Party:       m.Role{"buyer"},
			Token:       m.Ada,
			Deposits:    m.SetConstant("50000000"),
		},
	)
	assertJson(t, contract, `{"when":[{"case":{"into_account":{"role_token":"seller"},"party":{"role_token":"buyer"},"of_token":{"currency_symbol":"","token_name":""},"deposits":50000000},"then":"close"}],"timeout":1666078977926,"timeout_continuation":"close"}`)
}

func TestTypes_Notify(t *testing.T) {
	contract := setupWhenContract(
		m.Notify{
			If: m.ValueGT{
				Value: m.UseValue{"val"},
				Gt:    m.SetConstant("10"),
			},
		},
	)
	assertJson(t, contract, `{"when":[{"case":{"notify_if":{"value":{"use_value":"val"},"gt":10}},"then":"close"}],"timeout":1666078977926,"timeout_continuation":"close"}`)
}
