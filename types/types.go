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

type categories struct {
	categorySet map[string]empty
}

func (c categories) AddCategory(category string) {
	c.categorySet[category] = empty{}
}

func (c categories) HasCategory(category string) bool {
	_, exist := c.categorySet[category]
	return exist
}

func (c categories) Copy() categories {
	categorySet := map[string]empty{}
	for category := range c.categorySet {
		categorySet[category] = empty{}
	}
	return categories{categorySet: categorySet}
}

func makeCategories() categories {
	return categories{categorySet: map[string]empty{}}
}

type Loadable interface {
	LoadConfirm(Object) (Object, bool)
	Load(Object) Object
}

type Storable interface {
	Store(Object, Object)
}

type Environment interface {
	Loadable
	Storable
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
	var str string
	if b.inner {
		str = "true"
	} else {
		str = "false"
	}
	n, err := io.WriteString(w, str)
	return int64(n), err
}

func (b Boolean) Eval(env Environment) Object {
	return b
}

func MakeBoolean(b bool) Boolean {
	return Boolean{inner: b}
}

type Integer struct {
	categories
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
	return &Integer{categories: makeCategories(), inner: i}
}

type Float struct {
	categories
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
	return &Float{categories: makeCategories(), inner: f}
}

type String struct {
	categories
	inner string
}

func (s *String) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, s.inner)
	return int64(n), err
}

func (s *String) Eval(env Environment) Object {
	return s
}

func (s *String) LoadConfirm(key Object) (Object, bool) {
	var res Object
	exist := false
	switch casted := key.(type) {
	case *Integer:
		index := int(casted.inner)
		if 0 <= index && index < len(s.inner) {
			res = &String{
				categories: s.categories.Copy(),
				inner:      s.inner[index : index+1],
			}
			exist = true
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
							categories: s.categories.Copy(),
							inner:      s.inner[startInt:endInt],
						}
						exist = true
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
	return res, exist
}

func (s *String) Load(key Object) Object {
	res, _ := s.LoadConfirm(key)
	return res
}

func NewString(s string) *String {
	return &String{categories: makeCategories(), inner: s}
}

type Identifer String

func (i *Identifer) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, i.inner)
	return int64(n), err
}

func (i *Identifer) Eval(env Environment) Object {
	return env.Load((*String)(i))
}

func NewIdentifier(s string) *Identifer {
	return &Identifer{categories: makeCategories(), inner: s}
}

type Iterator interface {
	Next() (Object, bool)
}

type Iterable interface {
	Iter() Iterator
}

type List struct {
	categories
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

func (l *List) LoadConfirm(key Object) (Object, bool) {
	var res Object
	exist := false
	switch casted := key.(type) {
	case *Integer:
		index := int(casted.inner)
		if 0 <= index && index < len(l.inner) {
			res = l.inner[index]
			exist = true
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
							categories: l.categories.Copy(),
							inner:      l.inner[startInt:endInt],
						}
						exist = true
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
	return res, exist
}

func (l *List) Load(key Object) Object {
	res, _ := l.LoadConfirm(key)
	return res
}

func (l *List) Store(key, value Object) {
	integer, success := key.(*Integer)
	if success {
		index := int(integer.inner)
		if 0 <= index && index < len(l.inner) {
			l.inner[index] = value
		}
	}
}

type ListIterator struct {
	categories
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
	channel := make(chan Object)
	go func() {
		for _, value := range l.inner {
			channel <- value
		}
		close(channel)
	}()
	return &ListIterator{categories: l.categories, receive: channel}
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

type Appliable interface {
	Apply(Environment, *List) Object
}

func (l *List) Eval(env Environment) Object {
	var res Object
	if size := len(l.inner); size == 0 {
		res = None
	} else {
		value0 := l.inner[0].Eval(env)
		if f, success := value0.(Appliable); success {
			l2 := &List{
				categories: l.categories.Copy(),
				inner:      l.inner[1:],
			}
			res = f.Apply(env, l2)
		} else {
			l2 := &List{
				categories: l.categories.Copy(),
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
	return &List{categories: makeCategories()}
}

type BaseEnvironment struct {
	NoneType
	objects map[string]Object
}

func (b BaseEnvironment) LoadConfirm(key Object) (Object, bool) {
	var res Object
	exist := false
	str, success := key.(*String)
	if success {
		res, exist = b.objects[str.inner]
	} else {
		res = None
	}
	return res, exist
}

func (b BaseEnvironment) Load(key Object) Object {
	res, _ := b.LoadConfirm(key)
	return res
}

func (b BaseEnvironment) Store(key, value Object) {
	str, success := key.(*String)
	if success {
		b.objects[str.inner] = value
	}
}

func MakeBaseEnvironment() BaseEnvironment {
	return BaseEnvironment{objects: map[string]Object{}}
}

type LocalEnvironment struct {
	NoneType
	local  BaseEnvironment
	parent Environment
}

func (l *LocalEnvironment) LoadConfirm(key Object) (Object, bool) {
	res, exist := l.local.LoadConfirm(key)
	if !exist {
		res, exist = l.parent.LoadConfirm(key)
	}
	return res, exist
}

func (l *LocalEnvironment) Load(key Object) Object {
	res, _ := l.LoadConfirm(key)
	return res
}

func (l *LocalEnvironment) Store(key, value Object) {
	l.local.Store(key, value)
}

func NewLocalEnvironment(env Environment) *LocalEnvironment {
	return &LocalEnvironment{local: MakeBaseEnvironment(), parent: env}
}

type Native struct {
	NoneType
	inner func(Environment, *List) Object
}

func (n Native) Apply(env Environment, l *List) Object {
	return n.inner(env, l)
}

func makeNative(f func(Environment, *List) Object) Native {
	return Native{inner: f}
}
