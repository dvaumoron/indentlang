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
	"reflect"
)

type BaseEnvironment struct {
	NoneType
	objects map[string]Object
}

func (b BaseEnvironment) loadConfirm(key string) (Object, bool) {
	res, exist := b.objects[key]
	if !exist {
		res = None
	}
	return res, exist
}

func (b BaseEnvironment) Load(key Object) Object {
	str, success := key.(*String)
	var res Object
	if success {
		res, _ = b.loadConfirm(str.Inner)
	} else {
		res = None
	}
	return res
}

func (b BaseEnvironment) LoadStr(key string) Object {
	res, _ := b.loadConfirm(key)
	return res
}

func (b BaseEnvironment) Store(key, value Object) {
	str, success := key.(*String)
	if success {
		b.objects[str.Inner] = value
	}
}

func (b BaseEnvironment) StoreStr(key string, value Object) {
	b.objects[key] = value
}

func (b BaseEnvironment) Delete(key Object) {
	str, success := key.(*String)
	if success {
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
	NoneType
	local  BaseEnvironment
	parent Environment
}

func (l *LocalEnvironment) Load(key Object) Object {
	str, success := key.(*String)
	var res Object
	if success {
		res = l.LoadStr(str.Inner)
	} else {
		res = None
	}
	return res
}

func (l *LocalEnvironment) LoadStr(key string) Object {
	res, exist := l.local.loadConfirm(key)
	if !exist {
		res = l.parent.LoadStr(key)
	}
	return res
}

func (l *LocalEnvironment) Store(key, value Object) {
	l.local.Store(key, value)
}

func (l *LocalEnvironment) StoreStr(key string, value Object) {
	l.local.StoreStr(key, value)
}

func (l *LocalEnvironment) Delete(key Object) {
	l.local.Delete(key)
}

func (l *LocalEnvironment) DeleteStr(key string) {
	l.local.DeleteStr(key)
}

func (l *LocalEnvironment) CopyTo(other Environment) {
	l.local.CopyTo(other)
}

func NewLocalEnvironment(env Environment) *LocalEnvironment {
	return &LocalEnvironment{local: MakeBaseEnvironment(), parent: env}
}

type DataEnvironment struct {
	NoneType
	loadConfirm func(string) (Object, bool)
	parent      Environment
}

func (d *DataEnvironment) Load(key Object) Object {
	str, success := key.(*String)
	var res Object
	if success {
		res = d.LoadStr(str.Inner)
	} else {
		res = None
	}
	return res
}

func (d *DataEnvironment) LoadStr(key string) Object {
	res, exist := d.loadConfirm(key)
	if !exist {
		res = d.parent.LoadStr(key)
	}
	return res
}

func (d *DataEnvironment) Store(key, value Object) {
	d.parent.Store(key, value)
}

func (d *DataEnvironment) StoreStr(key string, value Object) {
	d.parent.StoreStr(key, value)
}

func (d *DataEnvironment) Delete(key Object) {
	d.parent.Delete(key)
}

func (d *DataEnvironment) DeleteStr(key string) {
	d.parent.DeleteStr(key)
}

func (d *DataEnvironment) CopyTo(other Environment) {
	d.parent.CopyTo(other)
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
	return &DataEnvironment{loadConfirm: loadConfirm, parent: env}
}
