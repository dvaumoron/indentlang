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

const returnName = "Return"

type returnMarker struct{}

type userAppliable struct {
	types.NoneType
	retrieveArgs func(types.Environment, *types.List) types.Environment
	body         *types.List
	evalReturn   func(types.Environment, types.Object) types.Object
}

func (u *userAppliable) Apply(env types.Environment, args *types.List) (res types.Object) {
	res = types.None
	local := u.retrieveArgs(env, args)
	defer func() {
		if r := recover(); r != nil {
			_, ok := r.(returnMarker)
			if ok {
				res = u.evalReturn(env, local.LoadStr(returnName))
			} else {
				panic(r)
			}
		}
	}()
	types.ForEach(u.body, func(line types.Object) bool {
		line.Eval(local)
		return true
	})
	return
}

func newUserFunction(declared types.Object, body *types.List) *userAppliable {
	var retrieveArgs func(types.Environment, *types.List) types.Environment
	switch casted := declared.(type) {
	case *types.Identifer:
		argName := casted.Inner
		retrieveArgs = func(env types.Environment, args *types.List) types.Environment {
			local := emptyFunctionArgs(env, args)
			local.StoreStr(argName, evalList(env, args))
			return local
		}
	case *types.List:
		if size := casted.SizeInt(); size == 0 {
			retrieveArgs = emptyFunctionArgs
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
				local := emptyFunctionArgs(env, args)
				it := args.Iter()
				for _, id := range ids {
					value, _ := it.Next()
					local.StoreStr(id, value.Eval(env))
				}
				return local
			}
		}
	default:
		retrieveArgs = emptyFunctionArgs
	}
	return &userAppliable{retrieveArgs: retrieveArgs, body: body, evalReturn: noEvalReturn}
}

func evalList(env types.Environment, args *types.List) *types.List {
	evaluated := types.NewList()
	types.ForEach(args, func(value types.Object) bool {
		evaluated.Add(value.Eval(env))
		return true
	})
	return evaluated
}

var functionReturnForm = types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
	var res types.Object = types.None
	if size := args.SizeInt(); size == 1 {
		res = args.Load(types.NewInteger(0)).Eval(env)
	} else if size > 1 {
		res = evalList(env, args)
	}
	env.StoreStr(returnName, res)
	panic(returnMarker{})
})

func emptyFunctionArgs(env types.Environment, args *types.List) types.Environment {
	local := types.NewLocalEnvironment(env)
	local.StoreStr(returnName, functionReturnForm)
	return local
}

func noEvalReturn(env types.Environment, res types.Object) types.Object {
	return res
}

func funcForm(env types.Environment, args *types.List) types.Object {
	if args.SizeInt() > 2 {
		it := args.Iter()
		funcName, _ := it.Next()
		declared, _ := it.Next()
		body := types.NewList()
		body.AddAll(it)
		env.Store(funcName, newUserFunction(declared, body))
	}
	return types.None
}

func lambdaForm(env types.Environment, args *types.List) (res types.Object) {
	if args.SizeInt() > 1 {
		it := args.Iter()
		declared, _ := it.Next()
		body := types.NewList()
		body.AddAll(it)
		res = newUserFunction(declared, body)
	} else {
		res = types.None
	}
	return
}

func callForm(env types.Environment, args *types.List) types.Object {
	// TODO
	return types.None
}

func newUserMacro(declared types.Object, body *types.List) *userAppliable {
	var retrieveArgs func(types.Environment, *types.List) types.Environment
	switch casted := declared.(type) {
	case *types.Identifer:
		argName := casted.Inner
		retrieveArgs = func(env types.Environment, args *types.List) types.Environment {
			local := emptyMacroArgs(env, args)
			local.StoreStr(argName, args)
			return local
		}
	case *types.List:
		if size := casted.SizeInt(); size == 0 {
			retrieveArgs = emptyMacroArgs
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
				local := emptyMacroArgs(env, args)
				it := args.Iter()
				for _, id := range ids {
					value, _ := it.Next()
					local.StoreStr(id, value)
				}
				return local
			}
		}
	default:
		retrieveArgs = emptyMacroArgs
	}
	return &userAppliable{retrieveArgs: retrieveArgs, body: body, evalReturn: doEvalReturn}
}

var macroReturnForm = types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
	env.StoreStr(returnName, args.Load(types.NewInteger(0)))
	panic(returnMarker{})
})

func emptyMacroArgs(env types.Environment, args *types.List) types.Environment {
	local := types.NewLocalEnvironment(env)
	local.StoreStr(returnName, macroReturnForm)
	return local
}

func doEvalReturn(env types.Environment, res types.Object) types.Object {
	return res.Eval(env)
}

func macroForm(env types.Environment, args *types.List) types.Object {
	if args.SizeInt() > 2 {
		it := args.Iter()
		macroName, _ := it.Next()
		declared, _ := it.Next()
		body := types.NewList()
		body.AddAll(it)
		env.Store(macroName, newUserMacro(declared, body))
	}
	return types.None
}
