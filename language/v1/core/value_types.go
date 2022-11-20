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
	"fmt"
	"math/big"
)

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
func (v TimeIntervalValue) isValue() {}
func (v UseValue) isValue()          {}
func (v Cond) isValue()              {}

// "Three of the Value terms look up information in the Marlowe state:
// AvailableMoney, ChoiceValue & UseValue" (§2.1.5)

// "AvailableMoney p t reports the amount of token t in the internal account of party p" (§2.1.5)
type AvailableMoney struct {
	Amount  Token     `json:"amount_of_token"`
	Account AccountId `json:"in_account"`
}

// "ChoiceValue i reports the most recent value chosen for choice i, or zero if
// no such choice has been made" (§2.1.5)
type ChoiceValue struct {
	Value ChoiceId `json:"value_of_choice"`
}

// "UseValue i reports the most recent value of the variable i, or zero
// if that variable has not yet been set to a value." (§2.1.5)
type UseValue struct {
	Value ValueId `json:"use_value"`
}

type TimeIntervalValue string

const TimeIntervalStart TimeIntervalValue = "time_interval_start"
const TimeIntervalEnd TimeIntervalValue = "time_interval_end"

// "Constant v evaluates to the integer v, while NegValue x, AddValue x y, SubValue x y,
// MulValue x y, and DivValue x y provide the common arithmetic operations -
// x, x + y, x − y, x ∗ y, and x / y, where division always rounds (truncates)
// its result towards zero." (§2.1.5)
type Constant Integer

// Make Constant a custom type for the JSON marshaller, converting big.Int to a string
func (i Constant) MarshalJSON() ([]byte, error) {
	i2 := big.Int(i)
	return []byte(fmt.Sprintf(`%s`, i2.String())), nil
}

func SetConstant(s string) Constant {
	bInt := big.NewInt(0)
	num, _ := bInt.SetString(s, 10)
	return Constant(*num)
}

type NegValue struct {
	Neg Value `json:"negate"`
}

// add() value (addition)
type AddValue struct {
	Add Value `json:"add"`
	To  Value `json:"and"`
}

// multiply() value (multiplication)
type MulValue struct {
	Multiply Value `json:"multiply"`
	By       Value `json:"times"`
}

// subtract() value (subtraction)
type SubValue struct {
	Subtract Value `json:"minus"`
	From     Value `json:"value"`
}

// div() value (division)
type DivValue struct {
	Divide Value `json:"divide"`
	By     Value `json:"by"`
}

// "Cond b x y represents a condition expression that evaluates to x if b is true
// and to y otherwise." (§2.1.5)
type Cond struct {
	Observation bool
	IfTrue      Value
	IfFalse     Value
}

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
	Both Observation `json:"both"`
	And  Observation `json:"and"`
}

func (o AndObs) isObservation() {}
func (o AndObs) isValue()       {}

type OrObs struct {
	Either Observation `json:"either"`
	Or     Observation `json:"or"`
}

func (o OrObs) isObservation() {}
func (o OrObs) isValue()       {}

type NotObs struct {
	Not Observation `json:"not"`
}

func (o NotObs) isObservation() {}
func (o NotObs) isValue()       {}

// "For the observations, the ChoseSomething i term reports whether a choice i
// has been made thus far in the contract." (§2.1.5)

type ChoseSomething struct {
	Choice ChoiceId `json:"chose_something_for"`
}

func (o ChoseSomething) isObservation() {}
func (o ChoseSomething) isValue()       {}

// "Value comparisons x < y, x ≤ y, x > y, x ≥ y, and x = y are represented
// by ValueLT x y, ValueLE x y, ValueGT x y, ValueGE x y, and ValueEQ x y." (§2.1.5)

type ValueGE struct {
	Value Value `json:"value"`
	Ge    Value `json:"ge_than"`
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
	Value Value `json:"value"`
	Le    Value `json:"le_than"`
}

func (o ValueLE) isObservation() {}
func (o ValueLE) isValue()       {}

type ValueEQ struct {
	Value Value `json:"value"`
	Eq    Value `json:"equal_to"`
}

func (o ValueEQ) isObservation() {}
func (o ValueEQ) isValue()       {}

// "The terms TrueObs and FalseObs provide the logical constants true and false." (§2.1.5)

type BoolObs bool

func (o BoolObs) isObservation() {}
func (o BoolObs) isValue()       {}

const TrueObs BoolObs = true
const FalseObs BoolObs = false
