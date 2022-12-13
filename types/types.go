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

import (
	"fmt"
	"io"
)

type Categorizable interface {
	AddCategory(string)
	HasCategory(string) bool
}

type empty = struct{}

type Categories struct {
	categorySet map[string]empty
}

func (c Categories) AddCategory(category string) {
	c.categorySet[category] = empty{}
}

func (c Categories) HasCategory(category string) bool {
	_, exist := c.categorySet[category]
	return exist
}

func (c Categories) Copy() Categories {
	categorySet := map[string]empty{}
	for category := range c.categorySet {
		categorySet[category] = empty{}
	}
	return Categories{categorySet: categorySet}
}

type Gettable interface {
	Get(Object) Object
}

type Settable interface {
	Set(Object, Object)
}

type Environment interface {
	Gettable
	Settable
}

type Object interface {
	Categorizable
	io.WriterTo
	Eval(Environment) Object
}

type NoneType empty

func (n NoneType) AddCategory(category string) {
}

func (n NoneType) HasCategory(category string) bool {
	return false
}

func (n NoneType) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (n NoneType) Eval(env Environment) Object {
	return None
}

var None = NoneType{}

type Boolean struct {
	NoneType
	inner bool
}

func (b Boolean) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, fmt.Sprint(b.inner))
	return int64(n), err
}

func (b Boolean) Eval(env Environment) Object {
	return b
}

func makeBoolean(b bool) Boolean {
	return Boolean{inner: b}
}

type Integer struct {
	Categories
	inner int64
}

func (i *Integer) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, fmt.Sprint(i.inner))
	return int64(n), err
}

func (i *Integer) Eval(env Environment) Object {
	return i
}

func NewInteger(i int64) *Integer {
	return &Integer{Categories: Categories{categorySet: map[string]empty{}}, inner: i}
}

type Float struct {
	Categories
	inner float64
}

func (f *Float) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, fmt.Sprint(f.inner))
	return int64(n), err
}

func (f *Float) Eval(env Environment) Object {
	return f
}

func NewFloat(f float64) *Float {
	return &Float{Categories: Categories{categorySet: map[string]empty{}}, inner: f}
}

type String struct {
	Categories
	inner string
}

func (s *String) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, s.inner)
	return int64(n), err
}

func (s *String) Eval(env Environment) Object {
	return s
}

func (s *String) Get(key Object) Object {
	var res Object
	switch casted := key.(type) {
	case *Integer:
		index := int(casted.inner)
		if 0 <= index && index < len(s.inner) {
			res = &String{
				Categories: s.Categories.Copy(),
				inner:      s.inner[index : index+1],
			}
		} else {
			res = None
		}
	case *List:
		if len(casted.inner) > 1 {
			start, success := casted.inner[0].(*Integer)
			if success {
				end, success := casted.inner[1].(*Integer)
				if success {
					startInt := int(start.inner)
					endInt := int(end.inner)
					if 0 <= startInt && startInt <= endInt && endInt < len(s.inner) {
						res = &String{
							Categories: s.Categories.Copy(),
							inner:      s.inner[startInt:endInt],
						}
					} else {
						res = None
					}
				} else {
					res = None
				}
			} else {
				res = None
			}
		} else {
			res = None
		}
	default:
		res = None
	}
	return res
}

func NewString(s string) *String {
	return &String{Categories: Categories{categorySet: map[string]empty{}}, inner: s}
}

type Identifer String

func (i *Identifer) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, i.inner)
	return int64(n), err
}

func (i *Identifer) Eval(env Environment) Object {
	return env.Get((*String)(i))
}

func NewIdentifier(s string) *Identifer {
	return &Identifer{Categories: Categories{categorySet: map[string]empty{}}, inner: s}
}

type Iterator interface {
	Next() (Object, bool)
}

type Iterable interface {
	Iter() Iterator
}

type List struct {
	Categories
	inner []Object
}

func (l *List) Add(value Object) {
	l.inner = append(l.inner, value)
}

func (l *List) AddAll(it Iterable) {
	it2 := it.Iter()
	exist := true
	var value Object
	for {
		value, exist = it2.Next()
		if !exist {
			break
		}
		l.Add(value)
	}
}

func (l *List) Get(key Object) Object {
	var res Object
	switch casted := key.(type) {
	case *Integer:
		index := int(casted.inner)
		if 0 <= index && index < len(l.inner) {
			res = l.inner[index]
		} else {
			res = None
		}
	case *List:
		if len(casted.inner) > 1 {
			start, success := casted.inner[0].(*Integer)
			if success {
				end, success := casted.inner[1].(*Integer)
				if success {
					startInt := int(start.inner)
					endInt := int(end.inner)
					if 0 <= startInt && startInt <= endInt && endInt < len(l.inner) {
						res = &List{
							Categories: l.Categories.Copy(),
							inner:      l.inner[startInt:endInt],
						}
					} else {
						res = None
					}
				} else {
					res = None
				}
			} else {
				res = None
			}
		} else {
			res = None
		}
	default:
		res = None
	}
	return res
}

func (l *List) Set(key, value Object) {
	integer, success := key.(*Integer)
	if success {
		index := int(integer.inner)
		if 0 <= index && index < len(l.inner) {
			l.inner[index] = value
		}
	}
}

type ListIterator struct {
	Categories
	receive <-chan Object
}

func (it *ListIterator) Iter() Iterator {
	return it
}

func (it *ListIterator) Next() (Object, bool) {
	value, exist := <-it.receive
	return value, exist
}

func (l *List) Iter() Iterator {
	canal := make(chan Object)
	go func() {
		for _, value := range l.inner {
			canal <- value
		}
		close(canal)
	}()
	return &ListIterator{Categories: l.Categories, receive: canal}
}

func WriteTo(it Iterable, w io.Writer) (int64, error) {
	var n int64
	var n2 int64
	var err error
	exist := true
	var value Object
	it2 := it.Iter()
	for err == nil {
		value, exist = it2.Next()
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

type Function interface {
	Apply(Environment, *List) Object
}

func (l *List) Eval(env Environment) Object {
	var res Object
	if size := len(l.inner); size == 0 {
		res = None
	} else {
		value0 := l.inner[0].Eval(env)
		if f, success := value0.(Function); success {
			res = f.Apply(env, l)
		} else {
			l2 := &List{
				Categories: Categories{categorySet: map[string]empty{}},
				inner:      make([]Object, 0, size),
			}
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
	return &List{Categories: Categories{categorySet: map[string]empty{}}}
}

type BaseEnvironment struct {
	NoneType
	objects map[string]Object
}

func (b BaseEnvironment) Get(key Object) Object {
	var res Object
	str, success := key.(*String)
	if success {
		res = b.objects[str.inner]
	} else {
		res = None
	}
	return res
}

func (b BaseEnvironment) Set(key, value Object) {
	str, success := key.(*String)
	if success {
		b.objects[str.inner] = value
	}
}

func makeBaseEnvironment() BaseEnvironment {
	return BaseEnvironment{objects: map[string]Object{}}
}

type LocalEnvironment struct {
	NoneType
	local  BaseEnvironment
	parent Environment
}

func (l *LocalEnvironment) Get(key Object) Object {
	res := l.local.Get(key)
	if res == nil {
		res = l.parent.Get(key)
	}
	return res
}

func (l *LocalEnvironment) Set(key, value Object) {
	l.Set(key, value)
}

func NewLocalEnvironment(env Environment) *LocalEnvironment {
	return &LocalEnvironment{local: makeBaseEnvironment(), parent: env}
}

type NativeFunction struct {
	NoneType
	inner func(Environment, *List) Object
}

func (n NativeFunction) Apply(env Environment, l *List) Object {
	return n.inner(env, l)
}

func makeNativeFunction(f func(Environment, *List) Object) NativeFunction {
	return NativeFunction{inner: f}
}
