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
	"time"
)

type BaseEnvironment struct {
	NoneType
	objects map[string]Object
}

func (b BaseEnvironment) LoadStr(key string) (Object, bool) {
	res, ok := b.objects[key]
	if !ok {
		return None, false
	}
	return res, true
}

func Load(env StringLoadable, key Object) Object {
	str, ok := key.(String)
	if !ok {
		return None
	}
	res, _ := env.LoadStr(string(str))
	return res
}

func (b BaseEnvironment) Load(key Object) Object {
	return Load(b, key)
}

func (b BaseEnvironment) Store(key, value Object) {
	str, ok := key.(String)
	if ok {
		b.objects[string(str)] = value
	}
}

func (b BaseEnvironment) StoreStr(key string, value Object) {
	b.objects[key] = value
}

func (b BaseEnvironment) Delete(key Object) {
	str, ok := key.(String)
	if ok {
		delete(b.objects, string(str))
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

func (b BaseEnvironment) Size() int {
	return len(b.objects)
}
func (b BaseEnvironment) Iter() Iterator {
	objectChannel := make(chan Object)
	it := &chanIterator{channel: objectChannel}
	go it.sendMapValue(b.objects)
	return it
}

type chanIterator struct {
	NoneType
	channel   chan Object
	cancelled bool
}

func (it *chanIterator) Iter() Iterator {
	return it
}

func (it *chanIterator) Next() (Object, bool) {
	value, ok := <-it.channel
	if !ok {
		return None, false
	}
	return value, true
}

func (it *chanIterator) Close() {
	it.cancelled = true
}

func (it *chanIterator) sendMapValue(objects map[string]Object) {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

ForLoop:
	for key, value := range objects {
		select {
		case it.channel <- NewList(String(key), value):
		case <-ticker.C:
			if it.cancelled {
				break ForLoop
			}
		}
	}
	close(it.channel)
}

func MakeBaseEnvironment() BaseEnvironment {
	return BaseEnvironment{objects: map[string]Object{}}
}

type LocalEnvironment struct {
	BaseEnvironment
	parent Environment
}

func (l LocalEnvironment) LoadStr(key string) (Object, bool) {
	res, ok := l.BaseEnvironment.LoadStr(key)
	if ok {
		return res, true
	}
	return l.parent.LoadStr(key)
}

func (l LocalEnvironment) Load(key Object) Object {
	return Load(l, key)
}

func MakeLocalEnvironment(env Environment) LocalEnvironment {
	return LocalEnvironment{BaseEnvironment: MakeBaseEnvironment(), parent: env}
}

type DataEnvironment struct {
	loadData ConvertString
	Environment
}

func (d DataEnvironment) LoadStr(key string) (Object, bool) {
	res, ok := d.loadData(key)
	if ok {
		return res, true
	}
	return d.Environment.LoadStr(key)
}

func (d DataEnvironment) Load(key Object) Object {
	return Load(d, key)
}

func MakeDataEnvironment(data any, env Environment) DataEnvironment {
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
	return DataEnvironment{loadData: loadConfirm, Environment: env}
}

// Read only : all non Load* methods ore no-op.
type MergeEnvironment struct {
	NoneType
	creationEnv Environment
	callEnv     Environment
}

func (m MergeEnvironment) LoadStr(key string) (Object, bool) {
	res, ok := m.creationEnv.LoadStr(key)
	if ok {
		return res, true
	}
	return m.callEnv.LoadStr(key)
}

func (m MergeEnvironment) Load(key Object) Object {
	return Load(m, key)
}

func (m MergeEnvironment) Store(key, value Object) {
}

func (m MergeEnvironment) StoreStr(key string, value Object) {
}

func (m MergeEnvironment) Delete(key Object) {
}

func (m MergeEnvironment) DeleteStr(key string) {
}

func (m MergeEnvironment) CopyTo(other Environment) {
}

func MakeMergeEnvironment(creationEnv, callEnv Environment) MergeEnvironment {
	return MergeEnvironment{creationEnv: creationEnv, callEnv: callEnv}
}
