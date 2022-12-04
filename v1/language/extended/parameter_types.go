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
	core "github.com/menabrealabs/marlowe/v1/language/core"
)

// These types are part of Marlowe Extended rather than Marlowe Core
// Marlowe extended has not been formally specified.
type TimeConstant core.POSIXTime
type TimeParam string

func (t TimeConstant) IsTimeout() {}
func (t TimeParam) IsTimeout()    {}

func (t TimeParam) ToCore(key string) {}

type ConstantParam string

func (c ConstantParam) IsValue() {}
