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

type Categorizable interface {
	AddCategory(string)
	HasCategory(string) bool
}

type Loadable interface {
	Load(Object) Object
}

type Storable interface {
	Loadable
	Store(Object, Object)
}

type StringLoadable interface {
	LoadStr(string) Object
}

type Environment interface {
	Storable
	Delete(Object)
	StringLoadable
	StoreStr(string, Object)
	DeleteStr(string)
	CopyTo(Environment)
}

type Object interface {
	Categorizable
	io.WriterTo
	Eval(Environment) Object
}

type Addable interface {
	Add(Object)
	AddAll(Iterable)
}

type Sizable interface {
	Size() *Integer
}

type Iterator interface {
	Iterable
	Next() (Object, bool)
}

type Iterable interface {
	Object
	Iter() Iterator
}

type Appliable interface {
	Apply(Environment, Iterable) Object
}
