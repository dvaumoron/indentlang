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
	categories
	inner []Object
}

func (l *List) Add(value Object) {
	l.inner = append(l.inner, value)
}

// If the action func return false that break the loop and ForEach return false too.
func ForEach(it Iterable, action func(Object) bool) bool {
	exist := true
	it2 := it.Iter()
	for exist {
		var value Object
		value, exist = it2.Next()
		if exist {
			exist = action(value)
		}
	}
	return exist
}

func (l *List) AddAll(it Iterable) {
	ForEach(it, func(value Object) bool {
		l.Add(value)
		return true
	})
}

func (l *List) Load(key Object) Object {
	var res Object
	switch casted := key.(type) {
	case *Integer:
		index := int(casted.Inner)
		if 0 <= index && index < len(l.inner) {
			res = l.inner[index]
		} else {
			res = None
		}
	case *List:
		if args := casted.inner; len(args) > 1 {
			start, success := args[0].(*Integer)
			if success {
				startInt := int(start.Inner)
				end, success := args[1].(*Integer)
				if success {
					endInt := int(end.Inner)
					if 0 <= startInt && startInt <= endInt && endInt < len(l.inner) {
						res = &List{
							categories: l.categories.Copy(),
							inner:      l.inner[startInt:endInt],
						}
					} else {
						res = None
					}
				} else {
					if 0 <= startInt && startInt < len(l.inner) {
						res = &List{
							categories: l.categories.Copy(),
							inner:      l.inner[startInt:],
						}
					} else {
						res = None
					}
				}
			} else {
				end, success := args[1].(*Integer)
				if success {
					endInt := int(end.Inner)
					if 0 <= endInt && endInt < len(l.inner) {
						res = &List{
							categories: l.categories.Copy(),
							inner:      l.inner[:endInt],
						}
					} else {
						res = None
					}
				} else {
					res = l
				}
			}
		} else {
			res = l
		}
	default:
		res = None
	}
	return res
}

func (l *List) Store(key, value Object) {
	integer, success := key.(*Integer)
	if success {
		index := int(integer.Inner)
		if 0 <= index && index < len(l.inner) {
			l.inner[index] = value
		}
	}
}

func (l *List) Size() *Integer {
	return NewInteger(int64(len(l.inner)))
}

func (l *List) SizeInt() int {
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
	value, exist := <-it.receiver
	if !exist {
		value = None
	}
	return value, exist
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

func WriteTo(it Iterable, w io.Writer) (int64, error) {
	var n, n2 int64
	var err error
	it2 := it.Iter()
	for err == nil {
		value, exist := it2.Next()
		if !exist {
			break
		}
		n2, err = value.WriteTo(w)
		n += n2
	}
	return n, err
}

func (l *List) WriteTo(w io.Writer) (int64, error) {
	return WriteTo(l, w)
}

func (l *List) Eval(env Environment) Object {
	it := l.Iter()
	res, exist := it.Next()
	if exist {
		value0 := res.Eval(env)
		if f, success := value0.(Appliable); success {
			res = f.Apply(env, it)
		} else {
			l2 := &List{categories: l.categories.Copy(), inner: make([]Object, 0, len(l.inner))}
			l2.Add(value0)
			for _, value := range l.inner[1:] {
				l2.Add(value.Eval(env))
			}
			res = l2
		}
	}
	return res
}

func NewList() *List {
	return &List{categories: makeCategories()}
}
