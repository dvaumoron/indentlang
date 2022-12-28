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
	var res types.Object = types.None
	ok := r.current < r.end
	if ok {
		r.current += r.step
		res = types.Integer(r.current)
	}
	return res, ok
}

func rangeFunc(env types.Environment, args types.Iterable) types.Object {
	var start int64
	var end int64
	step := int64(1)
	it := args.Iter()
	arg0, _ := it.Next()
	i0, ok := arg0.Eval(env).(types.Integer)
	if ok {
		arg1, _ := it.Next()
		var i1 types.Integer
		i1, ok = arg1.Eval(env).(types.Integer)
		if ok {
			start = int64(i0)
			end = int64(i1)
			arg2, _ := it.Next()
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
