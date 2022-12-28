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
	"fmt"
	"strconv"
	"strings"

	"github.com/dvaumoron/indentlang/types"
)

func boolConvFunc(env types.Environment, args types.Iterable) types.Object {
	arg, _ := args.Iter().Next()
	return types.Boolean(extractBoolean(arg.Eval(env)))
}

func extractBoolean(o types.Object) bool {
	res := true
	switch casted := o.(type) {
	case types.NoneType:
		res = false
	case types.Boolean:
		res = bool(casted)
	case types.Integer:
		res = casted != 0
	case types.Float:
		res = casted != 0
	case types.Sizable:
		res = casted.Size() != 0
	}
	return res
}

func intConvFunc(env types.Environment, args types.Iterable) types.Object {
	arg, _ := args.Iter().Next()
	return types.Integer(extractInteger(arg.Eval(env)))
}

func extractInteger(o types.Object) int64 {
	var res int64
	switch casted := o.(type) {
	case types.Boolean:
		if casted {
			res = 1
		}
	case types.Integer:
		res = int64(casted)
	case types.Float:
		res = int64(casted)
	case types.String:
		temp, err := strconv.ParseInt(string(casted), 10, 64)
		if err == nil {
			res = temp
		}
	}
	return res
}

func floatConvFunc(env types.Environment, args types.Iterable) types.Object {
	arg, _ := args.Iter().Next()
	return types.Float(extractFloat(arg.Eval(env)))
}

func extractFloat(o types.Object) float64 {
	var res float64
	switch casted := o.(type) {
	case types.Boolean:
		if casted {
			res = 1
		}
	case types.Integer:
		res = float64(casted)
	case types.Float:
		res = float64(casted)
	case types.String:
		temp, err := strconv.ParseFloat(string(casted), 64)
		if err == nil {
			res = temp
		}
	}
	return res
}

func stringConvFunc(env types.Environment, args types.Iterable) types.Object {
	arg, _ := args.Iter().Next()
	return types.String(extractString(arg.Eval(env)))
}

func extractString(o types.Object) string {
	res := ""
	switch casted := o.(type) {
	case types.Boolean:
		if casted {
			res = "true"
		} else {
			res = "false"
		}
	case types.Integer:
		res = fmt.Sprint(int64(casted))
	case types.Float:
		res = fmt.Sprint(float64(casted))
	case types.String:
		res = string(casted)
	case types.Iterable:
		it := casted.Iter()
		var builder strings.Builder
		builder.WriteRune('(')
		first, ok := it.Next()
		if ok {
			builder.WriteString(extractString(first))
			types.ForEach(it, func(el types.Object) bool {
				builder.WriteRune(' ')
				builder.WriteString(extractString(el))
				return true
			})
		}
		builder.WriteRune(')')
		res = builder.String()
	}
	return res
}

func listFunc(env types.Environment, args types.Iterable) types.Object {
	return types.NewList().AddAll(newEvalIterator(args, env))
}

func dictFunc(env types.Environment, args types.Iterable) types.Object {
	res := types.MakeBaseEnvironment()
	types.ForEach(args, func(arg types.Object) bool {
		it, ok := arg.Eval(env).(types.Iterable)
		if ok {
			it2 := it.Iter()
			var key types.Object
			key, ok = it2.Next()
			if ok {
				value, ok := it2.Next()
				if ok {
					res.Store(key, value)
				}
			}
		}
		return ok
	})
	return res
}
