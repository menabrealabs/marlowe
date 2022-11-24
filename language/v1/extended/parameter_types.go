package language

import (
	c "github.com/menabrealabs/marlowe/language/v1/core"
)

// These types are part of Marlowe Extended rather than Marlowe Core
// Marlowe extended has not been formally specified.
type TimeConstant c.POSIXTime
type TimeParam string

func (t TimeConstant) isTimeout() {}
func (t TimeParam) isTimeout()    {}

type ConstantParam string

func (c ConstantParam) isValue() {}
