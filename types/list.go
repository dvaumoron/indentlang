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
package types

import "io"

type List struct {
	categories map[string]NoneType
	inner      []Object
}

func (l *List) AddCategory(category string) {
	l.categories[category] = None
}

func (l *List) HasCategory(category string) bool {
	_, ok := l.categories[category]
	return ok
}

func (l *List) CopyCategories() map[string]NoneType {
	copy := map[string]NoneType{}
	for category := range l.categories {
		copy[category] = None
	}
	return copy
}

func (l *List) Add(value Object) {
	l.inner = append(l.inner, value)
}

// If the action func return false that break the loop.
func ForEach(it Iterable, action func(Object) bool) {
	ok := true
	it2 := it.Iter()
	for ok {
		var value Object
		value, ok = it2.Next()
		if ok {
			ok = action(value)
		}
	}
}

func (l *List) AddAll(it Iterable) *List {
	ForEach(it, func(value Object) bool {
		l.Add(value)
		return true
	})
	return l
}

func convertToInt(arg Object, init int) int {
	res := init
	casted, ok := arg.(Integer)
	if ok {
		res = int(casted)
	}
	return res
}

func extractIndex(args []Object, max int) (int, int) {
	var arg0, arg1 Object = None, None
	switch len(args) {
	case 0:
		// nothing to do, initialisation to None is enough
	default:
		// here, there is always enough element in the slice
		arg1 = args[1]
		fallthrough // allow initialisation of arg0
	case 1:
		// without fallthrough arg1 is None
		arg0 = args[0]
	}
	return convertToInt(arg0, 0), convertToInt(arg1, max)
}

func (l *List) LoadInt(index int) Object {
	var res Object = None
	if 0 <= index && index < len(l.inner) {
		res = l.inner[index]
	}
	return res
}

func (l *List) Load(key Object) Object {
	var res Object = None
	switch casted := key.(type) {
	case Integer:
		res = l.LoadInt(int(casted))
	case Float:
		res = l.LoadInt(int(casted))
	case *List:
		max := len(l.inner)
		start, end := extractIndex(casted.inner, max)
		if 0 <= start && start <= end && end <= max {
			res = &List{categories: l.CopyCategories(), inner: l.inner[start:end]}
		}
	}
	return res
}

func (l *List) Store(key, value Object) {
	integer, ok := key.(Integer)
	if ok {
		index := int(integer)
		if 0 <= index && index < len(l.inner) {
			l.inner[index] = value
		}
	}
}

func (l *List) Size() int {
	return len(l.inner)
}

type chanIterator struct {
	NoneType
	receiver <-chan Object
}

func (it *chanIterator) Iter() Iterator {
	return it
}

func (it *chanIterator) Next() (Object, bool) {
	value, ok := <-it.receiver
	if !ok {
		value = None
	}
	return value, ok
}

func (l *List) Iter() Iterator {
	channel := make(chan Object)
	go sendListValue(l.inner, channel)
	return &chanIterator{receiver: channel}
}

func sendListValue(objects []Object, transmitter chan<- Object) {
	for _, value := range objects {
		transmitter <- value
	}
	close(transmitter)
}

func (l *List) WriteTo(w io.Writer) (int64, error) {
	var n, n2 int64
	var err error
	for _, value := range l.inner {
		n2, err = value.WriteTo(w)
		n += n2
		if err != nil {
			break
		}
	}
	return n, err
}

func (l *List) Eval(env Environment) Object {
	it := l.Iter()
	res, ok := it.Next()
	if ok {
		value0 := res.Eval(env)
		if appliable, ok := value0.(Appliable); ok {
			res = appliable.Apply(env, it)
		} else {
			l2 := &List{categories: l.CopyCategories(), inner: make([]Object, 0, len(l.inner))}
			l2.Add(value0)
			for _, value := range l.inner[1:] {
				l2.Add(value.Eval(env))
			}
			res = l2
		}
	}
	return res
}

func NewList(objects ...Object) *List {
	return &List{categories: map[string]NoneType{}, inner: objects}
}
