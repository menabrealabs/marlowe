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
package language

import "math/big"

// "2.1.6 Actions and inputs
//
// Actions and Inputs are closely related. An Action can be added in a list of
// Cases §2.1.7 as a way to declare the possible external Inputs a Party can
// include in a Transaction at a certain time.
// The different types of actions are:
//
// datatype Action = Deposit AccountId Party Token Value
// | Choice ChoiceId Bound list
// | Notify Observation" (§2.1.6)

type Action interface{ isAction() }

// "A Deposit a p t v makes a deposit of #v Tokens t from Party p into account a." (§2.1.6)
type Deposit struct {
	AccountId AccountId
	Party     Party
	Token     Token
	Value     Value
}

func (a Deposit) isAction() {}

// "A choice Choice i bs is made for a particular choice identified by the ChoiceId
// §2.1.4 i with a list of inclusive bounds bs on the values that are acceptable.
// For example, [Bound 0 0 , Bound 3 5 ] offers the choice of one of 0, 3, 4 and
// 5." (§2.1.6)
type Choice struct {
	ChoiceId ChoiceId `json:"for_choice"`
	Bounds   []Bound  `json:"choose_between"`
}

func (a Choice) isAction() {}

// "Choices – of integers – are identified by ChoiceId which is defined with a
// canonical name and the Party who had made the choice." (§2.1.4)
type ChoiceId struct {
	Name  string `json:"choice_name"`
	Owner Party  `json:"choice_owner"`
}

// "Choices are Bounded. As an argument for the Choice action §2.1.6, we pass
// a list of Bounds that limit the integer that we can choose. The Bound data
// type is a tuple of integers that represents an inclusive lower and upper
// bound." (§2.1.4)
type Bound struct {
	Upper uint64 `json:"from"`
	Lower uint64 `json:"to"`
}

// "A notification can be triggered by anyone as long as the Observation evaluates
// to true. If multiple Notify are present in the Case list, the first one with a
// true observation is matched." (§2.1.6)
type Notify struct {
	Observation Observation
}

func (a Notify) isAction() {}

// "For each Action, there is a corresponding Input that can be included
// inside a Transaction:
//
// type-synonym ChosenNum = int
// datatype Input = IDeposit AccountId Party Token int
// | IChoice ChoiceId ChosenNum
// | INotify"
type ChosenNum int
type Input interface{ isInput() }

// "Deposit uses a Value while IDeposit has the int it was evaluated to
// with evalValue §2.2.10." (§2.1.6)
type IDeposit struct {
	AccountId AccountId
	Party     Party
	Token     Token
	Value     big.Int
}

func (i IDeposit) isInput() {}

// "Choice defines a list of valid Bounds while IChoice has the actual ChosenNum." (§2.1.6)
type IChoice struct {
	ChoiceId  ChoiceId
	ChosenNum ChosenNum
}

func (i IChoice) isInput() {}

// "Notify has an Observation while INotify does not have arguments, the
// Observation must evaluate to true inside the Transaction." (§2.1.6)
type INotify struct{}

func (i INotify) isInput() {}
