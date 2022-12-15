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

type NoneType struct{}

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

type categories struct {
	categorySet map[string]NoneType
}

func (c categories) AddCategory(category string) {
	c.categorySet[category] = None
}

func (c categories) HasCategory(category string) bool {
	_, exist := c.categorySet[category]
	return exist
}

func (c categories) Copy() categories {
	categorySet := map[string]NoneType{}
	for category := range c.categorySet {
		categorySet[category] = None
	}
	return categories{categorySet: categorySet}
}

func makeCategories() categories {
	return categories{categorySet: map[string]NoneType{}}
}

type Boolean struct {
	NoneType
	Inner bool
}

func (b Boolean) WriteTo(w io.Writer) (int64, error) {
	var str string
	if b.Inner {
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
	return Boolean{Inner: b}
}

type Integer struct {
	categories
	Inner int64
}

func (i *Integer) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, fmt.Sprint(i.Inner))
	return int64(n), err
}

func (i *Integer) Eval(env Environment) Object {
	return i
}

func NewInteger(i int64) *Integer {
	return &Integer{categories: makeCategories(), Inner: i}
}

type Float struct {
	categories
	Inner float64
}

func (f *Float) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, fmt.Sprint(f.Inner))
	return int64(n), err
}

func (f *Float) Eval(env Environment) Object {
	return f
}

func NewFloat(f float64) *Float {
	return &Float{categories: makeCategories(), Inner: f}
}

type String struct {
	categories
	Inner string
}

func (s *String) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, s.Inner)
	return int64(n), err
}

func (s *String) Eval(env Environment) Object {
	return s
}

func (s *String) Load(key Object) Object {
	var res Object
	switch casted := key.(type) {
	case *Integer:
		index := int(casted.Inner)
		if 0 <= index && index < len(s.Inner) {
			res = &String{
				categories: s.categories.Copy(),
				Inner:      s.Inner[index : index+1],
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
					startInt := int(start.Inner)
					endInt := int(end.Inner)
					if 0 <= startInt && startInt <= endInt && endInt < len(s.Inner) {
						res = &String{
							categories: s.categories.Copy(),
							Inner:      s.Inner[startInt:endInt],
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

func (s *String) Size() *Integer {
	return NewInteger(int64(len(s.Inner)))
}

func NewString(s string) *String {
	return &String{categories: makeCategories(), Inner: s}
}

type Identifer struct {
	String
}

func (i *Identifer) Eval(env Environment) Object {
	return env.Load((&i.String))
}

func NewIdentifier(s string) *Identifer {
	return &Identifer{String: *NewString(s)}
}

type NativeAppliable struct {
	NoneType
	inner func(Environment, *List) Object
}

func (n NativeAppliable) Apply(env Environment, l *List) Object {
	return n.inner(env, l)
}

func MakeNativeAppliable(f func(Environment, *List) Object) NativeAppliable {
	return NativeAppliable{inner: f}
}
