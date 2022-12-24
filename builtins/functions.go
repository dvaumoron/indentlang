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
	retrieveArgs(types.Environment, types.Environment, types.Iterable)
	defaultRetrieveArgs(types.Environment, types.Environment, types.Iterator)
	evalReturn(types.Environment, types.Environment) types.Object
}

type noArgsKind struct {
	returnForm types.NativeAppliable
	evalArgs   func(types.Iterable, types.Environment) types.Iterable
	evalIdArgs func(types.Environment, types.Environment, []string, types.Iterator)
	evalObject func(types.Object, types.Environment) types.Object
}

func (n *noArgsKind) initEnv(env types.Environment) types.Environment {
	local := types.NewLocalEnvironment(env)
	local.StoreStr(returnName, n.returnForm)
	return local
}

func (*noArgsKind) retrieveArgs(types.Environment, types.Environment, types.Iterable) {
}

func (*noArgsKind) defaultRetrieveArgs(types.Environment, types.Environment, types.Iterator) {
}

func (n *noArgsKind) evalReturn(env types.Environment, local types.Environment) types.Object {
	return n.evalObject(local.LoadStr(hiddenReturnName), env)
}

var functionKind = &noArgsKind{
	returnForm: types.MakeNativeAppliable(func(env types.Environment, args types.Iterable) types.Object {
		var res types.Object
		evaluated := evalIterable(args, env)
		switch size := evaluated.SizeInt(); size {
		case 0:
			res = types.None
		case 1:
			res = evaluated.Load(types.NewInteger(0))
		default:
			res = evaluated
		}
		env.StoreStr(hiddenReturnName, res)
		panic(returnMarker{})
	}),
	evalArgs: func(args types.Iterable, env types.Environment) types.Iterable {
		return evalIterable(args, env)
	},
	evalIdArgs: func(env types.Environment, local types.Environment, ids []string, it types.Iterator) {
		for _, id := range ids {
			value, _ := it.Next()
			local.StoreStr(id, value.Eval(env))
		}
	},
	evalObject: func(res types.Object, env types.Environment) types.Object {
		return res
	},
}
var macroKind = &noArgsKind{
	returnForm: types.MakeNativeAppliable(func(env types.Environment, args types.Iterable) types.Object {
		arg0, _ := args.Iter().Next()
		env.StoreStr(hiddenReturnName, arg0.Eval(env))
		panic(returnMarker{})
	}),
	evalArgs: func(args types.Iterable, env types.Environment) types.Iterable {
		return args
	},
	evalIdArgs: defaultRetrieveArgs,
	evalObject: func(res types.Object, env types.Environment) types.Object {
		return res.Eval(env)
	},
}

type classicKind struct {
	*noArgsKind
	ids []string
}

func (c *classicKind) retrieveArgs(env types.Environment, local types.Environment, args types.Iterable) {
	c.evalIdArgs(env, local, c.ids, args.Iter())

}

func (c *classicKind) defaultRetrieveArgs(env types.Environment, local types.Environment, it types.Iterator) {
	defaultRetrieveArgs(env, local, c.ids, it)
}

func defaultRetrieveArgs(env types.Environment, local types.Environment, ids []string, it types.Iterator) {
	for _, id := range ids {
		value, _ := it.Next()
		local.StoreStr(id, value)
	}
}

func newClassicKind(kind *noArgsKind, ids []string) *classicKind {
	return &classicKind{noArgsKind: kind, ids: ids}
}

type varArgsKind struct {
	*noArgsKind
	argsName string
}

func (v *varArgsKind) retrieveArgs(env types.Environment, local types.Environment, args types.Iterable) {
	local.StoreStr(v.argsName, v.evalArgs(args, env).Iter())
}

func (v *varArgsKind) defaultRetrieveArgs(env types.Environment, local types.Environment, it types.Iterator) {
	local.StoreStr(v.argsName, it)
}

func newVarArgsKind(kind *noArgsKind, name string) *varArgsKind {
	return &varArgsKind{noArgsKind: kind, argsName: name}
}

type userAppliable struct {
	types.NoneType
	body *types.List
	kindAppliable
}

func (u *userAppliable) Apply(env types.Environment, args types.Iterable) (res types.Object) {
	local := u.initEnv(env)
	u.retrieveArgs(env, local, args)
	defer func() {
		if r := recover(); r != nil {
			_, ok := r.(returnMarker)
			if ok {
				res = u.evalReturn(env, local)
			} else {
				panic(r)
			}
		}
	}()
	evalBody(u.body, local)
	return types.None
}

func (u *userAppliable) defaultApply(env types.Environment, it types.Iterator) (res types.Object) {
	local := u.initEnv(env)
	u.defaultRetrieveArgs(env, local, it)
	defer func() {
		if r := recover(); r != nil {
			_, ok := r.(returnMarker)
			if ok {
				res = u.evalReturn(env, local)
			} else {
				panic(r)
			}
		}
	}()
	evalBody(u.body, local)
	return types.None
}

func evalBody(body *types.List, local types.Environment) {
	types.ForEach(body, func(line types.Object) bool {
		line.Eval(local)
		return true
	})
}

func newUserAppliable(declared types.Object, body *types.List, baseKind *noArgsKind) *userAppliable {
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

func evalIterable(args types.Iterable, env types.Environment) *types.List {
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
		id2, ok := id.(*types.Identifer)
		if ok {
			ids = append(ids, id2.Inner)
		}
		return ok
	})
	return ids
}

func funcForm(env types.Environment, args types.Iterable) types.Object {
	return appliableForm(env, args, functionKind)
}

func macroForm(env types.Environment, args types.Iterable) types.Object {
	return appliableForm(env, args, macroKind)
}

func appliableForm(env types.Environment, args types.Iterable, kind *noArgsKind) types.Object {
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
	it := args.Iter()
	declared, _ := it.Next()
	body := types.NewList()
	body.AddAll(it)
	var res types.Object
	if body.SizeInt() != 0 {
		res = newUserAppliable(declared, body, functionKind)
	} else {
		res = types.None
	}
	return res
}

func callFunc(env types.Environment, args types.Iterable) types.Object {
	var res types.Object = types.None
	it := args.Iter()
	arg0, _ := it.Next()
	arg1, _ := it.Next()
	it2, ok := arg1.(types.Iterable)
	if ok {
		switch f := arg0.(type) {
		case *userAppliable:
			res = f.defaultApply(env, it2.Iter())
		case types.Appliable:
			res = f.Apply(env, it2)
		}
	}
	return res
}