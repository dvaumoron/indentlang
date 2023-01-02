/*
 *
 * Copyright 2022 indentlang authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package builtins

import (
	"github.com/dvaumoron/indentlang/parser"
	"github.com/dvaumoron/indentlang/types"
)

const sumName = "+"
const minusName = "-"
const productName = "*"
const divideName = "/"
const floorDivideName = "//"
const remainderName = "%"

type cumulCarac struct {
	init       int64
	cumulInt   func(int64, int64) int64
	cumulFloat func(float64, float64) float64
}

type number interface {
	int64 | float64
}

func addNumber[N number](a, b N) N {
	return a + b
}

func multNumber[N number](a, b N) N {
	return a * b
}

var sumCarac = &cumulCarac{
	init: 0, cumulInt: addNumber[int64], cumulFloat: addNumber[float64],
}
var productCarac = &cumulCarac{
	init: 1, cumulInt: multNumber[int64], cumulFloat: multNumber[float64],
}

func sumFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return cumulFunc(env, itArgs, sumCarac)
}

func productFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return cumulFunc(env, itArgs, productCarac)
}

func cumulFunc(env types.Environment, itArgs types.Iterator, carac *cumulCarac) types.Object {
	cumul := carac.init
	cumulF := float64(cumul)
	cumulInt := carac.cumulInt
	cumulFloat := carac.cumulFloat
	condition := true
	hasFloat := false
	types.ForEach(itArgs, func(arg types.Object) bool {
		switch casted := arg.Eval(env).(type) {
		case types.Integer:
			cumul = cumulInt(cumul, int64(casted))
		case types.Float:
			hasFloat = true
			cumulF = cumulFloat(cumulF, float64(casted))
		default:
			condition = false
		}
		return condition
	})
	var res types.Object = types.None
	if condition {
		if hasFloat {
			res = types.Float(cumulFloat(float64(cumul), cumulF))
		} else {
			res = types.Integer(cumul)
		}
	}
	return res
}

func minusFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, _ := itArgs.Next()
	var res types.Object = types.None
	switch casted := arg0.Eval(env).(type) {
	case types.Integer:
		switch casted2 := arg1.Eval(env).(type) {
		case types.Integer:
			res = types.Integer(casted - casted2)
		case types.Float:
			res = types.Float(float64(casted) - float64(casted2))
		}
	case types.Float:
		switch casted2 := arg1.Eval(env).(type) {
		case types.Integer:
			res = types.Float(float64(casted) - float64(casted2))
		case types.Float:
			res = types.Float(casted - casted2)
		}
	}
	return res
}

func divideFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, _ := itArgs.Next()
	var res types.Object = types.None
	switch casted := arg0.Eval(env).(type) {
	case types.Integer:
		res = divideObject(float64(casted), arg1.Eval(env))
	case types.Float:
		res = divideObject(float64(casted), arg1.Eval(env))
	}
	return res
}

func divideObject(a float64, b types.Object) types.Object {
	var res types.Object = types.None
	switch casted := b.(type) {
	case types.Integer:
		if casted != 0 {
			res = types.Float(a / float64(casted))
		}
	case types.Float:
		if casted != 0 {
			res = types.Float(a / float64(casted))
		}
	}
	return res
}

func floorDivideOperator(a, b int64) int64 {
	return a / b
}

func remainderOperator(a, b int64) int64 {
	return a % b
}

func floorDivideFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return intOperatorFunc(env, itArgs, floorDivideOperator)
}

func remainderFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return intOperatorFunc(env, itArgs, remainderOperator)
}

func intOperatorFunc(env types.Environment, itArgs types.Iterator, intOperator func(int64, int64) int64) types.Object {
	arg0, _ := itArgs.Next()
	arg1, _ := itArgs.Next()
	var res types.Object = types.None
	a, ok := arg0.Eval(env).(types.Integer)
	if ok {
		b, _ := arg1.Eval(env).(types.Integer)
		if b != 0 { // non integer and zero are treated the same way thanks to type assertion
			res = types.Integer(intOperator(int64(a), int64(b)))
		}
	}
	return res
}

func sumSetForm(env types.Environment, itArgs types.Iterator) types.Object {
	return inplaceOperatorForm(env, itArgs, sumName)
}

func minusSetForm(env types.Environment, itArgs types.Iterator) types.Object {
	return inplaceOperatorForm(env, itArgs, minusName)
}

func productSetForm(env types.Environment, itArgs types.Iterator) types.Object {
	return inplaceOperatorForm(env, itArgs, productName)
}

func divideSetForm(env types.Environment, itArgs types.Iterator) types.Object {
	return inplaceOperatorForm(env, itArgs, divideName)
}

func floorDivideSetForm(env types.Environment, itArgs types.Iterator) types.Object {
	return inplaceOperatorForm(env, itArgs, floorDivideName)
}

func remainderSetForm(env types.Environment, itArgs types.Iterator) types.Object {
	return inplaceOperatorForm(env, itArgs, remainderName)
}

func inplaceOperatorForm(env types.Environment, itArgs types.Iterator, opId types.Identifier) types.Object {
	arg0, _ := itArgs.Next()
	opInst := types.NewList(opId, arg0).AddAll(itArgs)
	return types.NewList(types.Identifier(parser.SetName), arg0, opInst)
}
