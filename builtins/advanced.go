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
	"github.com/dvaumoron/indentlang/parser"
	"github.com/dvaumoron/indentlang/types"
)

type quoteIterator struct {
	types.NoneType
	inner types.Iterator
	env   types.Environment
}

func (q *quoteIterator) Iter() types.Iterator {
	return q
}

func (q *quoteIterator) Next() (types.Object, bool) {
	value, ok := q.inner.Next()
	return evalUnquote(value, q.env), ok
}

func evalUnquote(object types.Object, env types.Environment) types.Object {
	list, ok := object.(*types.List)
	res := object
	if ok && list.Size() != 0 {
		resList := types.NewList()
		value := list.LoadInt(0)
		id, _ := value.(types.Identifier)
		if id == parser.UnquoteName { // non indentifier are ""
			res = list.LoadInt(1).Eval(env)
		} else {
			resList.Add(evalUnquote(value, env))
			resList.AddAll(newQuoteIterator(list.Iter(), env))
			res = resList
		}
	}
	return res
}

func newQuoteIterator(it types.Iterator, env types.Environment) *quoteIterator {
	return &quoteIterator{inner: it, env: env}
}

func quoteForm(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	return evalUnquote(arg, env)
}

// user can make (Lambda (x) (Return x))
func unquoteFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	return arg.Eval(env)
}

func delForm(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, ok := itArgs.Next()
	if ok {
		var dict types.Environment
		dict, ok = arg0.Eval(env).(types.Environment)
		if ok {
			dict.Delete(arg1.Eval(env))
		}
	} else {
		var id types.Identifier
		id, ok = arg0.(types.Identifier)
		if ok {
			env.DeleteStr(string(id))
		}
	}
	return types.None
}

func addCategoryFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, ok := itArgs.Next()
	if ok {
		var list *types.List
		list, ok = arg0.Eval(env).(*types.List)
		if ok {
			var str types.String
			str, ok = arg1.Eval(env).(types.String)
			if ok {
				list.AddCategory(string(str))
			}
		}
	}
	return types.None
}

func hasCategoryFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	arg1, ok := itArgs.Next()
	res := false
	if ok {
		var list *types.List
		list, ok = arg0.Eval(env).(*types.List)
		if ok {
			var str types.String
			str, ok = arg1.Eval(env).(types.String)
			if ok {
				res = list.HasCategory(string(str))
			}
		}
	}
	return types.Boolean(res)
}

func addCustomRuleFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	rule, ok := arg.Eval(env).(types.Appliable)
	if ok {
		parser.CustomRules.Add(rule)
	}
	return types.None
}

func parseWordFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	str, _ := arg.Eval(env).(types.String)
	list := types.NewList()
	if str != "" { // non string and empty string are treated the same way thanks to type assertion
		parser.HandleClassicWord(string(str), list)
	}
	return list.LoadInt(0) // None if the list is still empty
}
