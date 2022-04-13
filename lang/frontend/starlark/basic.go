// Copyright 2022 The MIDI Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package starlark

import (
	"math"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "math",
	Members: starlark.StringDict{
		"ceil":      starlark.NewBuiltin("ceil", ceil),
		"copysign":  newBinaryBuiltin("copysign", math.Copysign),
		"fabs":      newUnaryBuiltin("fabs", math.Abs),
		"floor":     starlark.NewBuiltin("floor", floor),
		"mod":       newBinaryBuiltin("round", math.Mod),
		"pow":       newBinaryBuiltin("pow", math.Pow),
		"remainder": newBinaryBuiltin("remainder", math.Remainder),
		"round":     newUnaryBuiltin("round", math.Round),

		"exp":  newUnaryBuiltin("exp", math.Exp),
		"sqrt": newUnaryBuiltin("sqrt", math.Sqrt),

		"acos":  newUnaryBuiltin("acos", math.Acos),
		"asin":  newUnaryBuiltin("asin", math.Asin),
		"atan":  newUnaryBuiltin("atan", math.Atan),
		"atan2": newBinaryBuiltin("atan2", math.Atan2),
		"cos":   newUnaryBuiltin("cos", math.Cos),
		"hypot": newBinaryBuiltin("hypot", math.Hypot),
		"sin":   newUnaryBuiltin("sin", math.Sin),
		"tan":   newUnaryBuiltin("tan", math.Tan),

		"degrees": newUnaryBuiltin("degrees", degrees),
		"radians": newUnaryBuiltin("radians", radians),

		"acosh": newUnaryBuiltin("acosh", math.Acosh),
		"asinh": newUnaryBuiltin("asinh", math.Asinh),
		"atanh": newUnaryBuiltin("atanh", math.Atanh),
		"cosh":  newUnaryBuiltin("cosh", math.Cosh),
		"sinh":  newUnaryBuiltin("sinh", math.Sinh),
		"tanh":  newUnaryBuiltin("tanh", math.Tanh),

		"log": starlark.NewBuiltin("log", log),

		"gamma": newUnaryBuiltin("gamma", math.Gamma),

		"e":  starlark.Float(math.E),
		"pi": starlark.Float(math.Pi),
	},
}
