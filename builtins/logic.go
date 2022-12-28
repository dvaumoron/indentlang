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
		res = extractBoolean(arg.Eval(env))
		return res
	})
	return types.Boolean(res)
}

func orFunc(env types.Environment, itArgs types.Iterator) types.Object {
	res := false
	types.ForEach(itArgs, func(arg types.Object) bool {
		res = extractBoolean(arg.Eval(env))
		return !res
	})
	return types.Boolean(res)
}
