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

func ifForm(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, _ := itArgs.Next()
	test := extractBoolean(arg0.Eval(env))
	var res types.Object
	if test {
		res = arg1.Eval(env)
	} else {
		arg2, _ := itArgs.Next()
		res = arg2.Eval(env)
	}
	return res
}

func forForm(env types.Environment, itArgs types.Iterator) types.Object {
	res := types.NewList()
	arg0, _ := itArgs.Next()
	arg1, _ := itArgs.Next()
	it, ok := arg1.Eval(env).(types.Iterable)
	if ok {
		bloc := types.NewList().AddAll(itArgs)
		if bloc.Size() != 0 {
			switch casted := arg0.(type) {
			case types.Identifier:

				id := string(casted)
				types.ForEach(it, func(value types.Object) bool {
					env.StoreStr(id, value)
					evalBloc(bloc, res, env)
					return true
				})
			case *types.List:
				ids := extractIds(casted)
				types.ForEach(it, func(value types.Object) bool {
					it2, ok := value.(types.Iterable)
					if ok {
						storeArgsInIds(ids, it2.Iter(), env)
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
	ids := make([]string, 0, l.Size())
	types.ForEach(l, func(value types.Object) bool {
		id, ok := value.(types.Identifier)
		if ok {
			ids = append(ids, string(id))
		}
		return ok
	})
	return ids
}

func storeArgsInIds(ids []string, itArgs types.Iterator, env types.Environment) {
	for _, id := range ids {
		arg, _ := itArgs.Next()
		env.StoreStr(id, arg)
	}
}

func whileForm(env types.Environment, itArgs types.Iterator) types.Object {
	res := types.NewList()
	arg0, _ := itArgs.Next()
	bloc := types.NewList().AddAll(itArgs)
	if bloc.Size() != 0 {
		for extractBoolean(arg0.Eval(env)) {
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

func setForm(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, ok := itArgs.Next()
	if ok {
		switch casted := arg0.(type) {
		case types.Identifier:
			env.StoreStr(string(casted), arg1.Eval(env))
		case *types.List:
			it, ok := arg1.Eval(env).(types.Iterable)
			if ok {
				it2 := it.Iter()
				types.ForEach(casted, func(value types.Object) bool {
					id, ok := value.(types.Identifier)
					if ok {
						value, _ := it2.Next()
						env.StoreStr(string(id), value)
					}
					return ok
				})
			}
		}
	}
	return types.None
}

func getForm(env types.Environment, itArgs types.Iterator) types.Object {
	res, ok := itArgs.Next()
	if ok {
		res = res.Eval(env)
		types.ForEach(itArgs, func(value types.Object) bool {
			var current types.StringLoadable
			current, ok = res.(types.StringLoadable)
			res = types.None
			if ok {
				var id types.Identifier
				id, ok = value.(types.Identifier)
				if ok {
					res, ok = current.LoadStr(string(id))
				}
			}
			return ok
		})
	}
	return res
}

func loadFunc(env types.Environment, itArgs types.Iterator) types.Object {
	res, ok := itArgs.Next()
	if ok {
		res = res.Eval(env)
		types.ForEach(itArgs, func(key types.Object) bool {
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

func storeFunc(env types.Environment, itArgs types.Iterator) types.Object {
	evaluated := types.NewList().AddAll(newEvalIterator(itArgs, env))
	if size := evaluated.Size() - 2; size > 0 {
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
