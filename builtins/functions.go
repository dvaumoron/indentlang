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

// user can not directly use this kind of id (# start comment)
const hiddenReturnName = "#return"

type returnMarker struct{}

type kindAppliable interface {
	initEnv(types.Environment) types.Environment
	retrieveArgs(types.Environment, types.Environment, types.Iterator)
	defaultRetrieveArgs(types.Environment, types.Environment, types.Iterator)
	evalReturn(types.Environment, types.Environment) types.Object
}

type noArgsKind struct {
	returnForm types.NativeAppliable
	evalArgs   func(types.Iterator, types.Environment) types.Iterator
	evalObject func(types.Object, types.Environment) types.Object
}

func (n noArgsKind) initEnv(env types.Environment) types.Environment {
	local := types.MakeLocalEnvironment(env)
	local.StoreStr(returnName, n.returnForm)
	return local
}

func (noArgsKind) retrieveArgs(types.Environment, types.Environment, types.Iterator) {
}

func (noArgsKind) defaultRetrieveArgs(types.Environment, types.Environment, types.Iterator) {
}

func (n noArgsKind) evalReturn(env types.Environment, local types.Environment) types.Object {
	returnValue, _ := local.LoadStr(hiddenReturnName)
	return n.evalObject(returnValue, env)
}

func noEval[T any](o T, env types.Environment) T {
	return o
}

var functionKind = noArgsKind{
	returnForm: types.MakeNativeAppliable(func(env types.Environment, itArgs types.Iterator) types.Object {
		evaluated := types.NewList().AddAll(makeEvalIterator(itArgs, env))
		switch size := evaluated.Size(); size {
		case 0:
			env.StoreStr(hiddenReturnName, types.None)
		case 1:
			env.StoreStr(hiddenReturnName, evaluated.LoadInt(0))
		default:
			env.StoreStr(hiddenReturnName, evaluated)
		}
		panic(returnMarker{})
	}),
	evalArgs: func(itArgs types.Iterator, env types.Environment) types.Iterator {
		return makeEvalIterator(itArgs, env)
	},
	evalObject: noEval[types.Object],
}
var macroKind = noArgsKind{
	returnForm: types.MakeNativeAppliable(func(env types.Environment, itArgs types.Iterator) types.Object {
		arg0, _ := itArgs.Next()
		env.StoreStr(hiddenReturnName, arg0.Eval(env))
		panic(returnMarker{})
	}),
	evalArgs: noEval[types.Iterator],
	evalObject: func(res types.Object, env types.Environment) types.Object {
		return res.Eval(env)
	},
}

type classicKind struct {
	noArgsKind
	ids []string
}

func (c classicKind) retrieveArgs(env types.Environment, local types.Environment, itArgs types.Iterator) {
	storeArgsInIds(c.ids, c.evalArgs(itArgs, env), local)
}

func (c classicKind) defaultRetrieveArgs(env types.Environment, local types.Environment, itArgs types.Iterator) {
	storeArgsInIds(c.ids, itArgs, local)
}

func makeClassicKind(kind noArgsKind, ids []string) classicKind {
	return classicKind{noArgsKind: kind, ids: ids}
}

type varArgsKind struct {
	noArgsKind
	argsName string
}

func (v varArgsKind) retrieveArgs(env types.Environment, local types.Environment, itArgs types.Iterator) {
	local.StoreStr(v.argsName, v.evalArgs(itArgs, env))
}

func (v varArgsKind) defaultRetrieveArgs(env types.Environment, local types.Environment, itArgs types.Iterator) {
	local.StoreStr(v.argsName, itArgs)
}

func makeVarArgsKind(kind noArgsKind, name string) varArgsKind {
	return varArgsKind{noArgsKind: kind, argsName: name}
}

type userAppliable struct {
	types.NoneType
	creationEnv types.Environment
	body        *types.List
	kindAppliable
}

func (u userAppliable) Apply(callEnv types.Environment, args types.Iterable) (res types.Object) {
	local := u.initEnv(u.creationEnv)
	itArgs := args.Iter()
	defer itArgs.Close()
	u.retrieveArgs(callEnv, local, itArgs)
	defer u.manageReturn(callEnv, local, &res)
	evalBody(u.body, local)
	return types.None
}

func (u userAppliable) manageReturn(callEnv types.Environment, local types.Environment, res *types.Object) {
	if r := recover(); r != nil {
		if _, ok := r.(returnMarker); !ok {
			panic(r)
		}
		*res = u.evalReturn(callEnv, local)
	}
}

func (u userAppliable) ApplyWithData(data any, callEnv types.Environment, args types.Iterable) (res types.Object) {
	local := u.initEnv(types.MakeMergeEnvironment(types.MakeDataEnvironment(data, u.creationEnv), callEnv))
	itArgs := args.Iter()
	defer itArgs.Close()
	u.retrieveArgs(callEnv, local, itArgs)
	defer u.manageReturn(callEnv, local, &res)
	evalBody(u.body, local)
	return types.None
}

func (u userAppliable) defaultApply(callEnv types.Environment, itArgs types.Iterator) (res types.Object) {
	local := u.initEnv(u.creationEnv)
	u.defaultRetrieveArgs(callEnv, local, itArgs)
	defer u.manageReturn(callEnv, local, &res)
	evalBody(u.body, local)
	return types.None
}

func evalBody(body *types.List, local types.Environment) {
	types.ForEach(body, func(line types.Object) bool {
		line.Eval(local)
		return true
	})
}

func makeUserAppliable(env types.Environment, declared types.Object, body *types.List, baseKind noArgsKind) userAppliable {
	var kind kindAppliable = baseKind
	switch casted := declared.(type) {
	case types.Identifier:
		kind = makeVarArgsKind(baseKind, string(casted))
	case *types.List:
		if casted.Size() != 0 {
			kind = makeClassicKind(baseKind, extractIds(casted))
		}
	}
	return userAppliable{creationEnv: env, body: body, kindAppliable: kind}
}

func funcForm(env types.Environment, itArgs types.Iterator) types.Object {
	return appliableForm(env, itArgs, functionKind)
}

func macroForm(env types.Environment, itArgs types.Iterator) types.Object {
	return appliableForm(env, itArgs, macroKind)
}

func appliableForm(env types.Environment, itArgs types.Iterator, kind noArgsKind) types.Object {
	arg0, _ := itArgs.Next()
	name, ok := arg0.(types.Identifier)
	if ok {
		declared, _ := itArgs.Next()
		body := types.NewList().AddAll(itArgs)
		if body.Size() != 0 {
			env.StoreStr(string(name), makeUserAppliable(env, declared, body, kind))
		}
	}
	return types.None
}

func lambdaForm(env types.Environment, itArgs types.Iterator) types.Object {
	declared, _ := itArgs.Next()
	body := types.NewList().AddAll(itArgs)
	if body.Size() == 0 {
		return types.None
	}
	return makeUserAppliable(env, declared, body, functionKind)
}

func callFunc(env types.Environment, itArgs types.Iterator) types.Object {
	var res types.Object = types.None
	arg0, _ := itArgs.Next()
	arg1, _ := itArgs.Next()
	it, ok := arg1.(types.Iterable)
	if ok {
		switch casted := arg0.(type) {
		case userAppliable:
			it2 := it.Iter()
			defer it2.Close()
			res = casted.defaultApply(env, it2)
		case types.Appliable:
			res = casted.Apply(env, it)
		}
	}
	return res
}
