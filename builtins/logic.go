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

import "github.com/dvaumoron/indentlang/types"

func notFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	return types.Boolean(!extractBoolean(arg.Eval(env)))
}

func andFunc(env types.Environment, itArgs types.Iterator) types.Object {
	res := true
	types.ForEach(itArgs, func(arg types.Object) bool {
		// change the variable that will be returned in the caller
		res = extractBoolean(arg.Eval(env))
		return res
	})
	return types.Boolean(res)
}

func orFunc(env types.Environment, itArgs types.Iterator) types.Object {
	res := false
	types.ForEach(itArgs, func(arg types.Object) bool {
		// change the variable that will be returned in the caller
		res = extractBoolean(arg.Eval(env))
		return !res
	})
	return types.Boolean(res)
}

func equalsFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, ok := itArgs.Next()
	return types.Boolean(ok && equals(arg0.Eval(env), arg1.Eval(env)))
}

func equals(value0, value1 types.Object) bool {
	res := false
	switch casted0 := value0.(type) {
	case types.NoneType:
		_, res = value1.(types.NoneType)
	case types.Boolean:
		casted1, ok := value1.(types.Boolean)
		res = ok && (casted0 == casted1)
	case types.Integer:
		switch casted1 := value1.(type) {
		case types.Integer:
			res = casted0 == casted1
		case types.Float:
			res = types.Float(casted0) == casted1
		}
	case types.Float:
		switch casted1 := value1.(type) {
		case types.Integer:
			res = casted0 == types.Float(casted1)
		case types.Float:
			res = casted0 == casted1
		}
	case types.String:
		casted1, ok := value1.(types.String)
		res = ok && (casted0 == casted1)
	}
	return res
}

func notEqualsFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, ok := itArgs.Next()
	return types.Boolean(ok && !equals(arg0.Eval(env), arg1.Eval(env)))
}

type comparator struct {
	compareInt    func(int64, int64) bool
	compareFloat  func(float64, float64) bool
	compareString func(string, string) bool
}

type ordered interface {
	number | string
}

func greaterThan[O ordered](value0, value1 O) bool {
	return value0 > value1
}

func greaterEqual[O ordered](value0, value1 O) bool {
	return value0 >= value1
}

func lessThan[O ordered](value0, value1 O) bool {
	return value0 < value1
}

func lessEqual[O ordered](value0, value1 O) bool {
	return value0 <= value1
}

var greaterThanComparator = &comparator{
	compareInt:    greaterThan[int64],
	compareFloat:  greaterThan[float64],
	compareString: greaterThan[string],
}
var greaterEqualComparator = &comparator{
	compareInt:    greaterEqual[int64],
	compareFloat:  greaterEqual[float64],
	compareString: greaterEqual[string],
}
var lessThanComparator = &comparator{
	compareInt:    lessThan[int64],
	compareFloat:  lessThan[float64],
	compareString: lessThan[string],
}
var lessEqualComparator = &comparator{
	compareInt:    lessEqual[int64],
	compareFloat:  lessEqual[float64],
	compareString: lessEqual[string],
}

func greaterThanFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return compareFunc(env, itArgs, greaterThanComparator)
}

func greaterEqualFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return compareFunc(env, itArgs, greaterEqualComparator)
}

func lessThanFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return compareFunc(env, itArgs, lessThanComparator)
}

func lessEqualFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return compareFunc(env, itArgs, lessEqualComparator)
}

func compareFunc(env types.Environment, itArgs types.Iterator, c *comparator) types.Object {
	arg0, res := itArgs.Next()
	if res {
		res = false
		previousValue := arg0.Eval(env)
		types.ForEach(itArgs, func(currentArg types.Object) bool {
			currentValue := currentArg.Eval(env)
			// change the variable that will be returned in the caller
			res = compare(previousValue, currentValue, c)
			previousValue = currentValue
			return res
		})
	}
	return types.Boolean(res)
}

func compare(value0 types.Object, value1 types.Object, c *comparator) bool {
	res := false
	switch casted0 := value0.(type) {
	case types.Integer:
		switch casted1 := value1.(type) {
		case types.Integer:
			res = c.compareInt(int64(casted0), int64(casted1))
		case types.Float:
			res = c.compareFloat(float64(casted0), float64(casted1))
		}
	case types.Float:
		switch casted1 := value1.(type) {
		case types.Integer:
			res = c.compareFloat(float64(casted0), float64(casted1))
		case types.Float:
			res = c.compareFloat(float64(casted0), float64(casted1))
		}
	case types.String:
		casted1, ok := value1.(types.String)
		res = ok && c.compareString(string(casted0), string(casted1))
	}
	return res
}
