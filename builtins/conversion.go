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
	"strconv"
	"strings"

	"github.com/dvaumoron/indentlang/types"
)

func identifierConvFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	s, ok := arg.Eval(env).(types.String)
	if !ok {
		return types.None
	}
	return types.Identifier(s)
}

func boolConvFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	return types.Boolean(extractBoolean(arg.Eval(env)))
}

func extractBoolean(o types.Object) bool {
	switch casted := o.(type) {
	case types.NoneType:
		return false
	case types.Boolean:
		return bool(casted)
	case types.Integer:
		return casted != 0
	case types.Float:
		return casted != 0
	case types.Sizable:
		return casted.Size() != 0
	}
	return true
}

func intConvFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	return types.Integer(extractInteger(arg.Eval(env)))
}

func extractInteger(o types.Object) int64 {
	switch casted := o.(type) {
	case types.Boolean:
		if casted {
			return 1
		}
	case types.Integer:
		return int64(casted)
	case types.Float:
		return int64(casted)
	case types.String:
		temp, err := strconv.ParseInt(string(casted), 10, 64)
		if err == nil {
			return temp
		}
	}
	return 0
}

func floatConvFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	return types.Float(extractFloat(arg.Eval(env)))
}

func extractFloat(o types.Object) float64 {
	switch casted := o.(type) {
	case types.Boolean:
		if casted {
			return 1
		}
	case types.Integer:
		return float64(casted)
	case types.Float:
		return float64(casted)
	case types.String:
		temp, err := strconv.ParseFloat(string(casted), 64)
		if err == nil {
			return temp
		}
	}
	return 0
}

func stringConvFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	return types.String(extractString(arg.Eval(env)))
}

func extractString(o types.Object) string {
	switch casted := o.(type) {
	case types.Boolean:
		if casted {
			return "true"
		}
		return "false"
	case types.Integer:
		return strconv.FormatInt(int64(casted), 10)
	case types.Float:
		return strconv.FormatFloat(float64(casted), 'g', -1, 64)
	case types.String:
		return string(casted)
	case types.Iterable:
		it := casted.Iter()
		defer it.Close()
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
		return builder.String()
	}
	return ""
}

func listFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return types.NewList().AddAll(makeEvalIterator(itArgs, env))
}

func dictFunc(env types.Environment, itArgs types.Iterator) types.Object {
	res := types.MakeBaseEnvironment()
	types.ForEach(itArgs, func(arg types.Object) bool {
		it, ok := arg.Eval(env).(types.Iterable)
		if !ok {
			return false
		}

		it2 := it.Iter()
		defer it2.Close()
		key, ok := it2.Next()
		if !ok {
			return false
		}

		value, ok := it2.Next()
		if ok {
			res.Store(key, value)
		}
		return ok
	})
	return res
}

func xmlTagFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	str, _ := arg.Eval(env).(types.String)
	if str == "" {
		return types.None
	}
	return createXmlTag(string(str))
}
