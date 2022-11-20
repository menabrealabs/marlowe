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

import (
	"math/big"
)

// use arbitrary-precision integers similar to Haskell's Integer primative
type Integer big.Int

// "We should separate the notions of participant, role, and address in a Marlowe
// contract. A participant (or Party) in the contract can be represented by
// either a fixed Address or a Role.
//
//	type-synonym RoleName = ByteString
//
//	datatype Party =
//		Address Address
//		| Role RoleName
type Party interface {
	isParty()
}

// "An address party is defined by a Blockhain specific Address §1.4 and it cannot
// be traded (it is fixed for the lifetime of a contract).
type Address string

// "A Role, on the other hand, allows the participation of the contract to be
// dynamic. Any user that can prove to have permission to act as RoleName
// is able to carry out the actions assigned §2.1.6, and redeem the payments
// issued to that role. The roles could be implemented as tokens that can be
// traded. By minting multiple tokens for a particular role, several people can
// be given permission to act on behalf of that role simultaneously, this allows
// for more complex use cases." (§2.1.1)
type Role struct {
	Name string `json:"role_token"`
}

func (r Role) isParty()    {}
func (p Address) isParty() {}

// "Inspired by Cardano’s Multi-Asset tokens, Marlowe also supports to transact with different assets.
// A Token consists of a CurrencySymbol that represents the monetary policy of the Token and a TokenName
// which allows to have multiple tokens with the same monetary policy.
//
// 	datatype Token = Token CurrencySymbol TokenName
//
// The Marlowe semantics treats both types as opaque ByteString." (§2.1.2)

// In order to allow Token to be used with AccountId in the Account type as a map index for Account
// we cannot use []byte for these values, since arrays are not Comparable types--a requirement for map
// index types. Convert to a []byte using the []byte(string) function to get the []byte representations
// if necessary.
type Token struct {
	Symbol string `json:"currency_symbol"`
	Name   string `json:"token_name"`
}

// Belongs in Cardano-specific implementation semantics
var Ada Token = Token{} // empty token defaults to $ADA

// "The Timeouts that prevent us from waiting forever for external Inputs are
// represented by the number of milliseconds from the Unix Epoch.
//
// type-synonym POSIXTime = int
// type-synonym Timeout = POSIXTime
//
// The TimeInterval that defines the validity of a transaction is a tuple of
// exclusive start and end time.
//
// type-synonym TimeInterval = POSIXTime × POSIXTime
// type POSIXtime int
// type Timeout POSIXtime" (§1.4)

type POSIXTime int
type Timeout POSIXTime

// Spec specifies a tuple, but Go doesn't have that datatype natively
type TimeInterval struct {
	// start is exclusive and end is inclusive
	start, end POSIXTime
}

// These types are part of Marlowe Extended rather than Marlowe Core
// Marlowe extended has not been formally specified.
type TimeConstant POSIXTime
type TimeParam string

func (t TimeConstant) isTimeout() {}
func (t TimeParam) isTimeout()    {}

type Payee struct {
	Party Party
}

type AccountId Party

// Go lacks tuples; Account implements an intermediate data structure
// that is not a type within the Marlowe Core spec.
type Account struct {
	AccountId AccountId
	Token     Token
}

func (a Account) isPayee() {}

type Accounts map[Account]uint64 // This is a type in the Marlowe Core specs.

// "The last Values, TimeIntervalStart and TimeIntervalEnd, evaluate respectively
// to the start or end of the validity interval for the Marlowe transaction." (§2.1.5)
// type TimeIntervalStart TimeInterval
// type TimeIntervalEnd TimeInterval

// "2.1.8 State and Environment

// The internal state of a Marlowe contract consists of the current balances in
// each party’s account, a record of the most recent value of each type of choice,
// a record of the most recent value of each variable, and the lower bound for the
// current time that is used to refine time intervals and ensure TimeIntervalStart
// never decreases. The data for accounts, choices, and bound values are stored
// as association lists.

// record State = accounts :: Accounts
// choices :: (ChoiceId × ChosenNum) list
// boundValues :: (ValueId × int) list
// minTime :: POSIXTime
type State struct {
	Accounts    Accounts
	BoundValues map[ValueId]uint64
	MinTime     POSIXTime
}

// The execution environment of a Marlowe contract simply consists of the
// (inclusive) time interval within which the transaction is occurring.

// record Environment = timeInterval :: TimeInterval"
type Environment struct {
	TimeInterval TimeInterval
}
