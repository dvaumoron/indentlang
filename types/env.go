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

func (b BaseEnvironment) LoadStr(key string) (Object, bool) {
	res, ok := b.objects[key]
	if !ok {
		res = None
	}
	return res, ok
}

func Load(env StringLoadable, key Object) Object {
	str, ok := key.(*String)
	var res Object = None
	if ok {
		res, _ = env.LoadStr(str.Inner)
	}
	return res
}

func (b BaseEnvironment) Load(key Object) Object {
	return Load(b, key)
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

func (l *LocalEnvironment) LoadStr(key string) (Object, bool) {
	res, ok := l.BaseEnvironment.LoadStr(key)
	if !ok {
		res, ok = l.parent.LoadStr(key)
	}
	return res, ok
}

func (l *LocalEnvironment) Load(key Object) Object {
	return Load(l, key)
}

func NewLocalEnvironment(env Environment) *LocalEnvironment {
	return &LocalEnvironment{BaseEnvironment: MakeBaseEnvironment(), parent: env}
}

type DataEnvironment struct {
	loadData func(string) (Object, bool)
	Environment
}

func (d *DataEnvironment) LoadStr(key string) (Object, bool) {
	res, ok := d.loadData(key)
	if !ok {
		res, ok = d.Environment.LoadStr(key)
	}
	return res, ok
}

func (d *DataEnvironment) Load(key Object) Object {
	return Load(d, key)
}

func NewDataEnvironment(data any, env Environment) *DataEnvironment {
	loadConfirm := neverConfirm
	dataValue, isNil := indirect(reflect.ValueOf(data))
	if !isNil {
		switch dataValue.Kind() {
		case reflect.Struct:
			loadConfirm = loadFromStruct(dataValue)
		case reflect.Map:
			loadConfirm = loadFromMap(dataValue)
		}
	}
	return &DataEnvironment{loadData: loadConfirm, Environment: env}
}

// Read only : all non Load* methods ore no-op.
type MergeEnvironment struct {
	creationEnv Environment
	callEnv     Environment
}

func (m *MergeEnvironment) LoadStr(key string) (Object, bool) {
	res, ok := m.creationEnv.LoadStr(key)
	if !ok {
		res, ok = m.callEnv.LoadStr(key)
	}
	return res, ok
}

func (m *MergeEnvironment) Load(key Object) Object {
	return Load(m, key)
}

func (m *MergeEnvironment) Store(key, value Object) {
}

func (m *MergeEnvironment) StoreStr(key string, value Object) {
}

func (m *MergeEnvironment) Delete(key Object) {
}

func (m *MergeEnvironment) DeleteStr(key string) {
}

func (m *MergeEnvironment) CopyTo(other Environment) {
}

func NewMergeEnvironment(creationEnv, callEnv Environment) *MergeEnvironment {
	return &MergeEnvironment{creationEnv: creationEnv, callEnv: callEnv}
}
