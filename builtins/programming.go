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

func ifForm(env types.Environment, args types.Iterable) types.Object {
	it := args.Iter()
	arg0, _ := it.Next()
	arg1, _ := it.Next()
	test := extractBoolean(arg0.Eval(env))
	var res types.Object
	if test {
		res = arg1.Eval(env)
	} else {
		arg2, _ := it.Next()
		res = arg2.Eval(env)
	}
	return res
}

func forForm(env types.Environment, args types.Iterable) types.Object {
	res := types.NewList()
	it := args.Iter()
	arg0, _ := it.Next()
	arg1, _ := it.Next()
	it2, ok := arg1.Eval(env).(types.Iterable)
	if ok {
		bloc := types.NewList()
		bloc.AddAll(it)
		switch casted := arg0.(type) {
		case types.Identifer:
			id := string(casted.String)
			types.ForEach(it2, func(value types.Object) bool {
				env.StoreStr(id, value)
				evalBloc(bloc, res, env)
				return true
			})
		case *types.List:
			if casted.SizeInt() != 0 {
				ids := extractIds(casted)
				types.ForEach(it2, func(value types.Object) bool {
					it3, ok := value.(types.Iterable)
					if ok {
						it4 := it3.Iter()
						for _, id := range ids {
							value2, _ := it4.Next()
							env.StoreStr(id, value2)
						}
						evalBloc(bloc, res, env)
					}
					return ok
				})
			}
		}
	}
	return res
}

func extractIds(l *types.List) []string {
	ids := make([]string, 0, l.SizeInt())
	types.ForEach(l, func(value types.Object) bool {
		id, ok := value.(types.Identifer)
		if ok {
			ids = append(ids, string(id.String))
		}
		return ok
	})
	return ids
}

func whileForm(env types.Environment, args types.Iterable) types.Object {
	res := types.NewList()
	it := args.Iter()
	arg0, _ := it.Next()
	bloc := types.NewList()
	bloc.AddAll(it)
	if bloc.SizeInt() != 0 {
		for {
			if !extractBoolean(arg0.Eval(env)) {
				break
			}
			evalBloc(bloc, res, env)
		}
	}
	return res
}

func evalBloc(bloc *types.List, res *types.List, env types.Environment) {
	types.ForEach(bloc, func(line types.Object) bool {
		evaluated := line.Eval(env)
		_, ok := evaluated.(types.NoneType)
		if !ok {
			res.Add(evaluated)
		}
		return true
	})
}

func setForm(env types.Environment, args types.Iterable) types.Object {
	it := args.Iter()
	arg0, _ := it.Next()
	arg1, ok := it.Next()
	if ok {
		switch casted := arg0.(type) {
		case types.Identifer:
			env.StoreStr(string(casted.String), arg1.Eval(env))
		case *types.List:
			it2, ok := arg1.Eval(env).(types.Iterable)
			if ok {
				it3 := it2.Iter()
				types.ForEach(casted, func(value types.Object) bool {
					id, ok := value.(types.Identifer)
					if ok {
						value, _ := it3.Next()
						env.StoreStr(string(id.String), value)
					}
					return ok
				})
			}
		}
	}
	return types.None
}

func getForm(env types.Environment, args types.Iterable) types.Object {
	it := args.Iter()
	res, ok := it.Next()
	if ok {
		res = res.Eval(env)
		types.ForEach(it, func(value types.Object) bool {
			current, ok := res.(types.StringLoadable)
			res = types.None
			if ok {
				var id types.Identifer
				id, ok = value.(types.Identifer)
				if ok {
					res, ok = current.LoadStr(string(id.String))
				}
			}
			return ok
		})
	}
	return res
}

func loadFunc(env types.Environment, args types.Iterable) types.Object {
	it := args.Iter()
	res, ok := it.Next()
	if ok {
		res = res.Eval(env)
		types.ForEach(it, func(key types.Object) bool {
			current, ok := res.(types.Loadable)
			if ok {
				res = current.Load(key.Eval(env))
			} else {
				res = types.None
			}
			return ok
		})
	}
	return res
}

func storeFunc(env types.Environment, args types.Iterable) types.Object {
	evaluated := evalIterable(args, env)
	if size := evaluated.SizeInt() - 2; size > 0 {
		it := evaluated.Iter()
		arg, _ := it.Next()
		ok := true
		for i := 1; i < size; i++ {
			key, _ := it.Next()
			var current types.Loadable
			current, ok = arg.(types.Loadable)
			if ok {
				arg = current.Load(key)
			} else {
				break
			}
		}
		if ok {
			current, ok := arg.(types.Storable)
			if ok {
				key, _ := it.Next()
				value, _ := it.Next()
				current.Store(key, value)
			}
		}
	}
	return types.None
}

func evalIterable(args types.Iterable, env types.Environment) *types.List {
	evaluated := types.NewList()
	types.ForEach(args, func(arg types.Object) bool {
		evaluated.Add(arg.Eval(env))
		return true
	})
	return evaluated
}
