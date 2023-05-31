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

func (l *List) ImportCategories(other *List) {
	for category := range other.categories {
		l.categories[category] = None
	}
}

func (l *List) Add(value Object) {
	l.inner = append(l.inner, value)
}

// If the action func return false that break the loop.
func ForEach(it Iterable, action func(Object) bool) {
	it2 := it.Iter()
	defer it2.Close()
	for ok := true; ok; {
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
	casted, ok := arg.(Integer)
	if !ok {
		return init
	}
	return int(casted)
}

func extractIndex(args []Object, max int) (int, int) {
	switch len(args) {
	case 0:
		return 0, max
	case 1:
		return convertToInt(args[0], 0), max
	}
	return convertToInt(args[0], 0), convertToInt(args[1], max)
}

func (l *List) LoadInt(index int) Object {
	if index < 0 || index >= len(l.inner) {
		return None
	}
	return l.inner[index]
}

func (l *List) Load(key Object) Object {
	switch casted := key.(type) {
	case Integer:
		return l.LoadInt(int(casted))
	case Float:
		return l.LoadInt(int(casted))
	case *List:
		max := len(l.inner)
		start, end := extractIndex(casted.inner, max)
		if 0 <= start && start <= end && end <= max {
			return &List{categories: l.CopyCategories(), inner: l.inner[start:end]}
		}
	}
	return None
}

func (l *List) Store(key Object, value Object) {
	integer, ok := key.(Integer)
	if ok {
		index := int(integer)
		if index >= 0 && index < len(l.inner) {
			l.inner[index] = value
		}
	}
}

func (l *List) Size() int {
	return len(l.inner)
}

type listIterator struct {
	NoneType
	list    *List
	current int
}

func (it *listIterator) Iter() Iterator {
	return it
}

func (it *listIterator) Next() (Object, bool) {
	inner := it.list.inner
	current := it.current
	if current >= len(inner) {
		return None, false
	}
	it.current++
	return inner[current], true
}

func (it *listIterator) Close() {
}

func (l *List) Iter() Iterator {
	return &listIterator{list: l}
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
	defer it.Close()
	elem0, ok := it.Next()
	if !ok {
		return None
	}
	value0 := elem0.Eval(env)
	if appliable, ok := value0.(Appliable); ok {
		return appliable.Apply(env, it)
	}
	l2 := &List{categories: l.CopyCategories(), inner: make([]Object, 0, len(l.inner))}
	l2.Add(value0)
	for _, value := range l.inner[1:] {
		l2.Add(value.Eval(env))
	}
	return l2
}

func NewList(objects ...Object) *List {
	return &List{categories: map[string]NoneType{}, inner: objects}
}
