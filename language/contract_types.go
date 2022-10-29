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

// "2.1.7 Contracts
//
// Marlowe is a continuation-based language, this means that a Contract can
// either be a Close or another construct that recursively has a Contract. Eventually,
// all contracts end up with a Close construct.
//
// Case and Contract are defined in a mutually recursive way as follows:
// datatype Case = Case Action Contract
// and Contract = Close
// | Pay AccountId Payee Token Value Contract
// | If Observation Contract Contract
// | When Case list Timeout Contract
// | Let ValueId Value Contract
// | Assert Observation Contract" (§2.1.6)
type CaseStmt interface{ isCase() }

type Contract interface {
	CaseStmt
	isContract()
}

type Case struct {
	Case     CaseStmt
	Action   Action
	Contract Contract
}

func (c Case) isCase() {}

// "Close is the simplest contract, when we evaluate it, the execution is completed
// and we generate Payments §?? for the assets in the internal accounts to their
// default owners." (§2.1.6)
type CloseContract string

const Close CloseContract = "close"

func (c CloseContract) isContract() {}
func (c CloseContract) isCase()     {}

// "The contract Pay a p t v c, generates a Payment from the internal account a
// to a payee §2.1.3 p of #v Tokens and then continues to contract c. Warnings
// will be generated if the value v is not positive, or if there is not enough in the
// account to make the payment in full. In the latter case, a partial payment
// (of the available amount) is made." (§2.1.6)
type Pay struct {
	AccountId AccountId `json:"from_account"`
	Payee     Payee     `json:"to"`
	Token     Token     `json:"token"`
	Pay       Value     `json:"pay"`
	Continue  Contract  `json:"then"`
}

func (c Pay) isContract() {}
func (c Pay) isCase()     {}

// "The contract If obs x y allows branching. We continue to branch x if the
// Observation obs evaluates to true, or to branch y otherwise." (§2.1.6)
type If struct {
	Observation Observation `json:"if"`
	Then        Contract    `json:"then"`
	Else        Contract    `json:"else"`
}

func (c If) isContract() {}
func (c If) isCase()     {}

// "When is the most complex constructor for contracts, with the form When cs t c.
// The list cs contains zero or more pairs of Actions and Contract continuations.
// When we do a computeTransaction §2.2.1, we follow the continuation
// associated to the first Action that matches the Input. If no action is matched
// it returns a ApplyAllNoMatchError. If a valid Transaction is computed with
// a TimeInterval with a start time bigger than the Timeout t, the contingency
// continuation c is evaluated. The explicit timeout mechanism is what allows
// Marlowe to avoid waiting forever for external inputs." (§2.1.6)
type When struct {
	Cases    []Case
	Timeout  Timeout
	Continue Contract
}

func (c When) isContract() {}
func (c When) isCase()     {}

// "A Let contract Let i v c allows a contract to record a value using an identifier
// i. In this case, the expression v is evaluated, and the result is stored with
// the name i. The contract then continues as c. As well as allowing us to
// use abbreviations, this mechanism also means that we can capture and save
// volatile values that might be changing with time, e.g. the current price of oil,
// or the current time, at a particular point in the execution of the contract, to
// be used later on in contract execution." (§2.1.6)
type Let struct {
	Name     ValueId  `json:"let"`
	Value    Value    `json:"be"`
	Continue Contract `json:"then"`
}

func (c Let) isContract() {}
func (c Let) isCase()     {}

// "An assertion contract Assert b c does not have any effect on the state of
// the contract, it immediately continues as c, but it issues a warning if the
// observation b evaluates to false. It can be used to ensure that a property
// holds in a given point of the contract, since static analysis will fail if any
// execution causes a warning. The Assert term might be removed from future
// on-chain versions of Marlowe." (§2.1.6)
type Assert struct {
	Observation Observation `json:"assert"`
	Continue    Contract    `json:"then"`
}

func (c Assert) isContract() {}
func (c Assert) isCase()     {}
