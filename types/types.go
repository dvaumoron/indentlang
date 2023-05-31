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

func (n NoneType) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (n NoneType) Eval(env Environment) Object {
	return None
}

var None = NoneType{}

type Boolean bool

func (b Boolean) WriteTo(w io.Writer) (int64, error) {
	var str string
	if b {
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

type Integer int64

func (i Integer) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, fmt.Sprint(int64(i)))
	return int64(n), err
}

func (i Integer) Eval(env Environment) Object {
	return i
}

type Float float64

func (f Float) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, fmt.Sprint(float64(f)))
	return int64(n), err
}

func (f Float) Eval(env Environment) Object {
	return f
}

type String string

func (s String) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, string(s))
	return int64(n), err
}

func (s String) Eval(env Environment) Object {
	return s
}

func (s String) LoadInt(index int) Object {
	if index < 0 || index >= len(s) {
		return None
	}
	return s[index : index+1]
}

func (s String) Load(key Object) Object {
	switch casted := key.(type) {
	case Integer:
		return s.LoadInt(int(casted))
	case Float:
		return s.LoadInt(int(casted))
	case *List:
		max := len(s)
		start, end := extractIndex(casted.inner, max)
		if 0 <= start && start <= end && end <= max {
			return s[start:end]
		}
	}
	return None
}

func (s String) Size() int {
	return len(s)
}

type Identifier string

func (i Identifier) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (i Identifier) Eval(env Environment) Object {
	value, _ := env.LoadStr(string(i))
	return value
}

type NativeAppliable struct {
	NoneType
	inner func(Environment, Iterator) Object
}

func (n NativeAppliable) Apply(env Environment, it Iterable) Object {
	it2 := it.Iter()
	defer it2.Close()
	return n.inner(env, it2)
}

func (n NativeAppliable) ApplyWithData(data any, env Environment, it Iterable) Object {
	it2 := it.Iter()
	defer it2.Close()
	return n.inner(MakeDataEnvironment(data, env), it2)
}

func MakeNativeAppliable(f func(Environment, Iterator) Object) NativeAppliable {
	return NativeAppliable{inner: f}
}
