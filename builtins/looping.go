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

type rangeIterator struct {
	types.NoneType
	current int64
	end     int64
	step    int64
}

func (r *rangeIterator) Iter() types.Iterator {
	return r
}

func (r *rangeIterator) Next() (types.Object, bool) {
	var res types.Object = types.Integer(r.current)
	ok := r.current < r.end
	if ok {
		r.current += r.step
	}
	return res, ok
}

func (r *rangeIterator) Close() {
}

func rangeFunc(env types.Environment, itArgs types.Iterator) types.Object {
	var start int64
	var end int64
	step := int64(1)
	arg0, _ := itArgs.Next()
	i0, ok := arg0.Eval(env).(types.Integer)
	if ok {
		arg1, _ := itArgs.Next()
		var i1 types.Integer
		i1, ok = arg1.Eval(env).(types.Integer)
		if ok {
			start = int64(i0)
			end = int64(i1)
			arg2, _ := itArgs.Next()
			var i2 types.Integer
			i2, ok = arg2.Eval(env).(types.Integer)
			if ok {
				step = int64(i2)
			}
		} else {
			end = int64(i0)
		}
	}
	return &rangeIterator{current: start, end: end, step: step}
}

type enumerateIterator struct {
	types.NoneType
	inner types.Iterator
	count int64
}

func (e *enumerateIterator) Iter() types.Iterator {
	return e
}

func (e *enumerateIterator) Next() (types.Object, bool) {
	value, ok := e.inner.Next()
	count := e.count
	e.count++
	return types.NewList(types.Integer(count), value), ok
}

func (e *enumerateIterator) Close() {
	e.inner.Close()
}

func enumerateFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg, _ := itArgs.Next()
	it, ok := arg.Eval(env).(types.Iterable)
	if !ok {
		return types.None
	}
	return &enumerateIterator{inner: it.Iter()}
}

func iterFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	it, ok := arg0.Eval(env).(types.Iterable)
	if !ok {
		return types.None
	}
	return it.Iter()
}

func nextFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	it, ok := arg0.Eval(env).(types.Iterator)
	var res0 types.Object = types.None
	res1 := false
	if ok {
		res0, res1 = it.Next()
	}
	return types.NewList(res0, types.Boolean(res1))
}

func closeFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	it, ok := arg0.Eval(env).(types.Iterator)
	if ok {
		it.Close()
	}
	return types.None
}

func sizeFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	it, ok := arg0.Eval(env).(types.Sizable)
	if !ok {
		return types.None
	}
	return types.Integer(it.Size())
}

type evalIterator struct {
	types.NoneType
	inner types.Iterator
	env   types.Environment
}

func (e evalIterator) Iter() types.Iterator {
	return e
}

func (e evalIterator) Next() (types.Object, bool) {
	value, ok := e.inner.Next()
	return value.Eval(e.env), ok
}

func (e evalIterator) Close() {
	e.inner.Close()
}

func makeEvalIterator(it types.Iterator, env types.Environment) evalIterator {
	return evalIterator{inner: it, env: env}
}

func addFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	list, ok := arg0.Eval(env).(*types.List)
	if ok {
		list.AddAll(makeEvalIterator(itArgs, env))
	}
	return types.None
}

func addAllFunc(env types.Environment, itArgs types.Iterator) types.Object {
	arg0, _ := itArgs.Next()
	list, ok := arg0.Eval(env).(*types.List)
	if ok {
		types.ForEach(itArgs, func(arg types.Object) bool {
			it2, ok := arg.Eval(env).(types.Iterable)
			if ok {
				list.AddAll(it2)
			}
			return ok
		})
	}
	return types.None
}
