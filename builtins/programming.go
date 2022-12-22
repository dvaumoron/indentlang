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

func ifForm(env types.Environment, args *types.List) types.Object {
	var res types.Object = types.None
	if args.SizeInt() > 1 {
		it := args.Iter()
		arg0, _ := it.Next()
		test, ok := arg0.Eval(env).(types.Boolean)
		if ok {
			arg1, _ := it.Next()
			if test.Inner {
				res = arg1.Eval(env)
			} else {
				arg2, exist := it.Next()
				if exist {
					res = arg2.Eval(env)
				}
			}
		}

	}
	return res
}

func forForm(env types.Environment, args *types.List) types.Object {
	res := types.NewList()
	if args.SizeInt() > 2 {
		it := args.Iter()
		arg0, _ := it.Next()
		arg1, _ := it.Next()
		it2, success := arg1.Eval(env).(types.Iterable)
		if success {
			bloc := types.NewList()
			bloc.AddAll(it)
			switch casted := arg0.(type) {
			case *types.Identifer:
				id := casted.Inner
				types.ForEach(it2, func(value types.Object) bool {
					env.StoreStr(id, value)
					types.ForEach(bloc, func(line types.Object) bool {
						res.Add(line.Eval(env))
						return true
					})
					return true
				})
			case *types.List:
				ids := make([]string, 0, casted.SizeInt())
				types.ForEach(casted, func(id types.Object) bool {
					id2, success := id.(*types.Identifer)
					if success {
						ids = append(ids, id2.Inner)
					}
					return true
				})
				types.ForEach(it2, func(value types.Object) bool {
					it3, success := value.(types.Iterable)
					if success {
						it4 := it3.Iter()
						for _, id := range ids {
							value2, _ := it4.Next()
							env.StoreStr(id, value2)
						}
						types.ForEach(bloc, func(line types.Object) bool {
							res.Add(line.Eval(env))
							return true
						})
					}
					return success
				})
			}
		}
	}
	return res
}

func setForm(env types.Environment, args *types.List) types.Object {
	if args.SizeInt() > 1 {
		it := args.Iter()
		arg0, _ := it.Next()
		arg1, _ := it.Next()
		switch casted := arg0.(type) {
		case *types.Identifer:
			env.StoreStr(casted.Inner, arg1.Eval(env))
		case *types.List:
			it2, success := arg1.Eval(env).(types.Iterable)
			if success {
				it3 := it2.Iter()
				types.ForEach(casted, func(id types.Object) bool {
					id2, success := id.(*types.Identifer)
					if success {
						value, _ := it3.Next()
						env.StoreStr(id2.Inner, value)
					}
					return success
				})
			}
		}
	}
	return types.None
}

func getForm(env types.Environment, args *types.List) types.Object {
	it := args.Iter()
	res, ok := it.Next()
	if ok {
		types.ForEach(it, func(value types.Object) bool {
			current, ok := res.(types.StringLoadable)
			if ok {
				var id *types.Identifer
				id, ok = value.(*types.Identifer)
				if ok {
					res = current.LoadStr(id.Inner)
				}
			}
			return ok
		})
	}
	return res
}

func indexForm(env types.Environment, args *types.List) types.Object {
	it := args.Iter()
	res, ok := it.Next()
	if ok {
		types.ForEach(it, func(value types.Object) bool {
			current, ok := res.(types.Loadable)
			if ok {
				res = current.Load(value.Eval(env))
			}
			return ok
		})
	}
	return res
}
