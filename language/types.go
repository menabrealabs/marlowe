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

// use arbitrary-precision integers similar to Haskell's Integer primative
type IntegerString string

func (v IntegerString) isValue() {}

func (v IntegerString) Int() *big.Int {
	i := big.NewInt(0)
	i.SetString(string(v), 10)
	return i
}

// "We should separate the notions of participant, role, and address in a Marlowe
// contract. A participant (or Party) in the contract can be represented by
// either a fixed Address or a Role.
//
//	type-synonym RoleName = ByteString
//
//	datatype Party =
//		Address Address
//		| Role RoleName
//
// "An address party is defined by a Blockhain specific Address §1.4 and it cannot
// be traded (it is fixed for the lifetime of a contract).
type Party interface {
	isParty()
}

// "A Role, on the other hand, allows the participation of the contract to be
// dynamic. Any user that can prove to have permission to act as RoleName
// is able to carry out the actions assigned §2.1.6, and redeem the payments
// issued to that role. The roles could be implemented as tokens that can be
// traded. By minting multiple tokens for a particular role, several people can
// be given permission to act on behalf of that role simultaneously, this allows
// for more complex use cases." (§2.1.1)
type Role struct {
	RoleName string `json:"role_token"`
}

type Address string

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
	CurrencySymbol string `json:"currency_symbol"`
	TokenName      string `json:"token_name"`
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

// "Choices – of integers – are identified by ChoiceId which is defined with a
// canonical name and the Party who had made the choice." (§2.1.4)
type ChoiceId struct {
	ChoiceName  string
	ChoiceOwner Party
}

// "Choices are Bounded. As an argument for the Choice action §2.1.6, we pass
// a list of Bounds that limit the integer that we can choose. The Bound data
// type is a tuple of integers that represents an inclusive lower and upper
// bound." (§2.1.4)
type Bound struct {
	Upper, Lower uint64
}

// "We can store a Value in the Marlowe State §2.1.8 using the Let construct
// §2.1.7, and we use a ValueId to referrence it" (§2.1.5)
type ValueId string

// "Values and Observations are language terms that interact with most of the
// other constructs. Value evaluates to an integer and Observation evaluates to
// a boolean using evalValue §2.2.10 and evalObservation §2.2.11 respectively.
// They are defined in a mutually recursive way as follows:
//
// datatype Value = AvailableMoney AccountId Token
// | Constant int
// | NegValue Value
// | AddValue Value Value
// | SubValue Value Value
// | MulValue Value Value
// | DivValue Value Value
// | ChoiceValue ChoiceId
// | TimeIntervalStart
// | TimeIntervalEnd
// | UseValue ValueId
// | Cond Observation Value Value"
type Value interface{ isValue() }

func (v AvailableMoney) isValue()    {}
func (v Constant) isValue()          {}
func (v NegValue) isValue()          {}
func (v AddValue) isValue()          {}
func (v SubValue) isValue()          {}
func (v MulValue) isValue()          {}
func (v DivValue) isValue()          {}
func (v ChoiceValue) isValue()       {}
func (v TimeIntervalStart) isValue() {}
func (v TimeIntervalEnd) isValue()   {}
func (v UseValue) isValue()          {}
func (v Cond) isValue()              {}

// "Three of the Value terms look up information in the Marlowe state:
// AvailableMoney, ChoiceValue & UseValue" (§2.1.5)

// "AvailableMoney p t reports the amount of token t in the internal account of party p" (§2.1.5)
type AvailableMoney struct {
	Amount  Token
	Account AccountId
}

// "ChoiceValue i reports the most recent value chosen for choice i, or zero if
// no such choice has been made" (§2.1.5)
type ChoiceValue ChoiceId

// "UseValue i reports the most recent value of the variable i, or zero
// if that variable has not yet been set to a value." (§2.1.5)
type UseValue ValueId

// "Constant v evaluates to the integer v, while NegValue x, AddValue x y, SubValue x y,
// MulValue x y, and DivValue x y provide the common arithmetic operations -
// x, x + y, x − y, x ∗ y, and x / y, where division always rounds (truncates)
// its result towards zero." (§2.1.5)
type Constant uint64

type NegValue struct{ Value Value }

// add() value (addition)
type AddValue struct {
	Lhs Value
	Rhs Value
}

// multiply() value (multiplication)
type MulValue struct {
	Lhs Value
	Rhs Value
}

// subtract() value (subtraction)
type SubValue struct {
	Lhs Value
	Rhs Value
}

// div() value (division)
type DivValue struct {
	Lhs Value
	Rhs Value
}

// "Cond b x y represents a condition expression that evaluates to x if b is true
// and to y otherwise." (§2.1.5)
type Cond struct {
	Observation bool
	IfTrue      Value
	IfFalse     Value
}

// "The last Values, TimeIntervalStart and TimeIntervalEnd, evaluate respectively
// to the start or end of the validity interval for the Marlowe transaction." (§2.1.5)
type TimeIntervalStart TimeInterval
type TimeIntervalEnd TimeInterval

// "and Observation = AndObs Observation Observation
// | OrObs Observation Observation
// | NotObs Observation
// | ChoseSomething ChoiceId
// | ValueGE Value Value
// | ValueGT Value Value
// | ValueLT Value Value
// | ValueLE Value Value
// | ValueEQ Value Value
// | TrueObs
// | FalseObs" (§2.1.5)

type Observation interface {
	Value
	isObservation()
}

// "The logical operators ¬ x, x ∧ y, and x ∨ y are represented by the terms
// NotObs x, AndObs x y, and OrObs x y, respectively." (§2.1.5)

type AndObs struct {
	Lhs, Rhs Observation
}

func (o AndObs) isObservation() {}
func (o AndObs) isValue()       {}

type OrObs struct {
	Lhs, Rhs Observation
}

func (o OrObs) isObservation() {}
func (o OrObs) isValue()       {}

type NotObs struct {
	Observation Observation
}

func (o NotObs) isObservation() {}
func (o NotObs) isValue()       {}

// "For the observations, the ChoseSomething i term reports whether a choice i
// has been made thus far in the contract." (§2.1.5)

type ChoseSomething struct {
	Choice ChoiceId
}

func (o ChoseSomething) isObservation() {}
func (o ChoseSomething) isValue()       {}

// "Value comparisons x < y, x ≤ y, x > y, x ≥ y, and x = y are represented
// by ValueLT x y, ValueLE x y, ValueGT x y, ValueGE x y, and ValueEQ x y." (§2.1.5)

type ValueGE struct {
	Lhs, Rhs Value
}

func (o ValueGE) isObservation() {}
func (o ValueGE) isValue()       {}

type ValueGT struct {
	Value Value `json:"value"`
	Gt    Value `json:"gt"`
}

func (o ValueGT) isObservation() {}
func (o ValueGT) isValue()       {}

type ValueLT struct {
	Value Value `json:"value"`
	Lt    Value `json:"lt"`
}

func (o ValueLT) isObservation() {}
func (o ValueLT) isValue()       {}

type ValueLE struct {
	Lhs, Rhs Value
}

func (o ValueLE) isObservation() {}
func (o ValueLE) isValue()       {}

type ValueEQ struct {
	Lhs, Rhs Value
}

func (o ValueEQ) isObservation() {}
func (o ValueEQ) isValue()       {}

// "The terms TrueObs and FalseObs provide the logical constants true and false." (§2.1.5)

type TrueObs struct{}

func (o TrueObs) isObservation() {}
func (o TrueObs) isValue()       {}

type FalseObs struct{}

func (o FalseObs) isObservation() {}
func (o FalseObs) isValue()       {}

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
	ChoiceId ChoiceId
	Bounds   []Bound
}

func (a Choice) isAction() {}

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
	Value     uint64
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
	If   Observation `json:"if"`
	Then Contract    `json:"then"`
	Else Contract    `json:"else"`
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
	ValueId  ValueId  `json:"let"`
	Value    Value    `json:"be"`
	Contract Contract `json:"then"`
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
	Assert   Observation `json:"assert"`
	Continue Contract    `json:"then"`
}

func (c Assert) isContract() {}
func (c Assert) isCase()     {}

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
