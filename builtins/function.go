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

type userFunction struct {
	types.NoneType
	retrieveArgs func(types.Environment, *types.List) types.Environment
	body         *types.List
}

func (u *userFunction) Apply(env types.Environment, args *types.List) types.Object {
	local := u.retrieveArgs(env, args)
	res := types.NewList()
	types.ForEach(u.body, func(line types.Object) bool {
		res.Add(line.Eval(local))
		return true
	})
	return res
}

func newUserFunction(declared types.Object, body *types.List) *userFunction {
	var retrieveArgs func(types.Environment, *types.List) types.Environment
	switch casted := declared.(type) {
	case *types.Identifer:
		argName := casted.Inner
		retrieveArgs = func(env types.Environment, args *types.List) types.Environment {
			local := types.NewLocalEnvironment(env)
			evaluated := types.NewList()
			types.ForEach(args, func(value types.Object) bool {
				evaluated.Add(value.Eval(env))
				return true
			})
			local.StoreStr(argName, evaluated)
			return local
		}
	case *types.List:
		if size := casted.SizeInt(); size == 0 {
			retrieveArgs = emptyArgs
		} else {
			ids := make([]string, 0, size)
			types.ForEach(casted, func(id types.Object) bool {
				id2, success := id.(*types.Identifer)
				if success {
					ids = append(ids, id2.Inner)
				}
				return success
			})
			retrieveArgs = func(env types.Environment, args *types.List) types.Environment {
				local := types.NewLocalEnvironment(env)
				it := args.Iter()
				for _, id := range ids {
					value, _ := it.Next()
					local.StoreStr(id, value.Eval(env))
				}
				return local
			}
		}
	default:
		retrieveArgs = emptyArgs
	}
	return &userFunction{retrieveArgs: retrieveArgs, body: body}
}

func emptyArgs(env types.Environment, args *types.List) types.Environment {
	return types.NewLocalEnvironment(env)
}

func funcForm(env types.Environment, args *types.List) types.Object {
	if args.SizeInt() > 3 {
		it := args.Iter()
		key, _ := it.Next()
		declared, _ := it.Next()
		body := types.NewList()
		body.AddAll(it)
		env.Store(key, newUserFunction(declared, body))
	}
	return types.None
}

func lambdaForm(env types.Environment, args *types.List) (res types.Object) {
	if args.SizeInt() > 2 {
		it := args.Iter()
		declared, _ := it.Next()
		body := types.NewList()
		body.AddAll(it)
		res = newUserFunction(declared, body)
	} else {
		res = types.None
	}
	return res
}
