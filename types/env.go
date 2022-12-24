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

import "reflect"

type BaseEnvironment struct {
	NoneType
	objects map[string]Object
}

func (b BaseEnvironment) loadConfirm(key string) (Object, bool) {
	res, ok := b.objects[key]
	if !ok {
		res = None
	}
	return res, ok
}

func (b BaseEnvironment) Load(key Object) Object {
	str, ok := key.(*String)
	var res Object = None
	if ok {
		res, _ = b.loadConfirm(str.Inner)
	}
	return res
}

func (b BaseEnvironment) LoadStr(key string) Object {
	res, _ := b.loadConfirm(key)
	return res
}

func (b BaseEnvironment) Store(key, value Object) {
	str, ok := key.(*String)
	if ok {
		b.objects[str.Inner] = value
	}
}

func (b BaseEnvironment) StoreStr(key string, value Object) {
	b.objects[key] = value
}

func (b BaseEnvironment) Delete(key Object) {
	str, ok := key.(*String)
	if ok {
		delete(b.objects, str.Inner)
	}
}

func (b BaseEnvironment) DeleteStr(key string) {
	delete(b.objects, key)
}

func (b BaseEnvironment) CopyTo(other Environment) {
	for key, value := range b.objects {
		other.StoreStr(key, value)
	}
}

func MakeBaseEnvironment() BaseEnvironment {
	return BaseEnvironment{objects: map[string]Object{}}
}

type LocalEnvironment struct {
	BaseEnvironment
	parent Environment
}

func (l *LocalEnvironment) Load(key Object) Object {
	str, ok := key.(*String)
	var res Object = None
	if ok {
		res = l.LoadStr(str.Inner)
	}
	return res
}

func (l *LocalEnvironment) LoadStr(key string) Object {
	res, ok := l.loadConfirm(key)
	if !ok {
		res = l.parent.LoadStr(key)
	}
	return res
}

func NewLocalEnvironment(env Environment) *LocalEnvironment {
	return &LocalEnvironment{BaseEnvironment: MakeBaseEnvironment(), parent: env}
}

type DataEnvironment struct {
	loadConfirm func(string) (Object, bool)
	Environment
}

func (d *DataEnvironment) Load(key Object) Object {
	str, ok := key.(*String)
	var res Object = None
	if ok {
		res = d.LoadStr(str.Inner)
	}
	return res
}

func (d *DataEnvironment) LoadStr(key string) Object {
	res, ok := d.loadConfirm(key)
	if !ok {
		res = d.Environment.LoadStr(key)
	}
	return res
}

func NewDataEnvironment(data any, env Environment) *DataEnvironment {
	var loadConfirm func(string) (Object, bool)
	dataValue, isNil := indirect(reflect.ValueOf(data))
	if isNil {
		loadConfirm = neverConfirm
	} else {
		switch dataValue.Kind() {
		case reflect.Struct:
			loadConfirm = loadFromStruct(dataValue)
		case reflect.Map:
			if stringType.AssignableTo(dataValue.Type().Key()) {
				loadConfirm = loadFromMap(dataValue)
			} else {
				loadConfirm = neverConfirm
			}
		default:
			loadConfirm = neverConfirm
		}
	}
	return &DataEnvironment{loadConfirm: loadConfirm, Environment: env}
}
