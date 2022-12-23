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

type kindAppliable interface {
	initEnv(types.Environment) types.Environment
	retrieveArgs(types.Environment, types.Environment, types.Iterable)
	defaultRetrieveArgs(types.Environment, types.Environment, types.Iterator)
	evalReturn(types.Environment, types.Object) types.Object
}

type noArgsKind bool

var macroReturnForm = types.MakeNativeAppliable(func(env types.Environment, args types.Iterable) types.Object {
	arg0, _ := args.Iter().Next()
	env.StoreStr(returnName, arg0.Eval(env))
	panic(returnMarker{})
})

// Not considered as a function because it panic
var functionReturnForm = types.MakeNativeAppliable(func(env types.Environment, args types.Iterable) types.Object {
	var res types.Object
	evaluated := evalIterable(env, args)
	switch size := evaluated.SizeInt(); size {
	case 0:
		res = types.None
	case 1:
		res = evaluated.Load(types.NewInteger(0))
	default:
		res = evaluated
	}
	env.StoreStr(returnName, res)
	panic(returnMarker{})
})

func (b noArgsKind) initEnv(env types.Environment) types.Environment {
	local := types.NewLocalEnvironment(env)
	if b {
		local.StoreStr(returnName, macroReturnForm)
	} else {
		local.StoreStr(returnName, functionReturnForm)
	}
	return local
}

func (b noArgsKind) retrieveArgs(types.Environment, types.Environment, types.Iterable) {
}

func (b noArgsKind) defaultRetrieveArgs(types.Environment, types.Environment, types.Iterator) {
}

func (b noArgsKind) evalReturn(env types.Environment, res types.Object) types.Object {
	if b {
		res = res.Eval(env)
	}
	return res
}

var functionKind = noArgsKind(false)
var macroKind = noArgsKind(true)

type classicKind struct {
	noArgsKind
	ids []string
}

func (c *classicKind) retrieveArgs(env types.Environment, local types.Environment, args types.Iterable) {
	it := args.Iter()
	if c.noArgsKind {
		c.defaultRetrieveArgs(env, local, it)
	} else {
		for _, id := range c.ids {
			value, _ := it.Next()
			local.StoreStr(id, value.Eval(env))
		}
	}
}

func (c *classicKind) defaultRetrieveArgs(env types.Environment, local types.Environment, it types.Iterator) {
	for _, id := range c.ids {
		value, _ := it.Next()
		local.StoreStr(id, value)
	}
}

func newClassicKind(kind noArgsKind, ids []string) *classicKind {
	return &classicKind{noArgsKind: kind, ids: ids}
}

type varArgsKind struct {
	noArgsKind
	argsName string
}

func (v *varArgsKind) retrieveArgs(env types.Environment, local types.Environment, args types.Iterable) {
	if v.noArgsKind {
		local.StoreStr(v.argsName, args.Iter())
	} else {
		local.StoreStr(v.argsName, evalIterable(env, args).Iter())
	}
}

func (v *varArgsKind) defaultRetrieveArgs(env types.Environment, local types.Environment, it types.Iterator) {
	local.StoreStr(v.argsName, it)
}

func newVarArgsKind(kind noArgsKind, name string) *varArgsKind {
	return &varArgsKind{noArgsKind: kind, argsName: name}
}

type userAppliable struct {
	types.NoneType
	body *types.List
	kindAppliable
}

func (u *userAppliable) Apply(env types.Environment, args types.Iterable) (res types.Object) {
	res = types.None
	local := u.initEnv(env)
	u.retrieveArgs(env, local, args)
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

func (u *userAppliable) defaultApply(env types.Environment, it types.Iterator) (res types.Object) {
	res = types.None
	local := u.initEnv(env)
	u.defaultRetrieveArgs(env, local, it)
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

func newUserAppliable(declared types.Object, body *types.List, baseKind noArgsKind) *userAppliable {
	var kind kindAppliable = baseKind
	switch casted := declared.(type) {
	case *types.Identifer:
		kind = newVarArgsKind(baseKind, casted.Inner)
	case *types.List:
		if casted.SizeInt() != 0 {
			kind = newClassicKind(baseKind, extractIds(casted))
		}
	}
	return &userAppliable{body: body, kindAppliable: kind}
}

func evalIterable(env types.Environment, args types.Iterable) *types.List {
	evaluated := types.NewList()
	types.ForEach(args, func(value types.Object) bool {
		evaluated.Add(value.Eval(env))
		return true
	})
	return evaluated
}

func extractIds(l *types.List) []string {
	ids := make([]string, 0, l.SizeInt())
	types.ForEach(l, func(id types.Object) bool {
		id2, success := id.(*types.Identifer)
		if success {
			ids = append(ids, id2.Inner)
		}
		return success
	})
	return ids
}

func funcForm(env types.Environment, args types.Iterable) types.Object {
	return appliableForm(env, args, functionKind)
}

func macroForm(env types.Environment, args types.Iterable) types.Object {
	return appliableForm(env, args, macroKind)
}

func appliableForm(env types.Environment, args types.Iterable, kind noArgsKind) types.Object {
	it := args.Iter()
	arg0, _ := it.Next()
	name, ok := arg0.(*types.Identifer)
	if ok {
		declared, _ := it.Next()
		body := types.NewList()
		body.AddAll(it)
		if body.SizeInt() != 0 {
			env.StoreStr(name.Inner, newUserAppliable(declared, body, kind))
		}
	}
	return types.None
}

func lambdaForm(env types.Environment, args types.Iterable) types.Object {
	var res types.Object = types.None
	it := args.Iter()
	declared, _ := it.Next()
	body := types.NewList()
	body.AddAll(it)
	if body.SizeInt() != 0 {
		res = newUserAppliable(declared, body, functionKind)
	}
	return res
}

func callForm(env types.Environment, args types.Iterable) types.Object {
	var res types.Object = types.None
	it := args.Iter()
	arg0, _ := it.Next()
	arg1, _ := it.Next()
	it2, ok := arg1.(types.Iterable)
	if ok {
		switch f := arg0.(type) {
		case *userAppliable:
			f.defaultApply(env, it2.Iter())
		case types.Appliable:
			f.Apply(env, it2)
		}
	}
	return res
}
