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
	"strconv"
)

type ConvertString func(string) (Object, bool)
type ExtractIterator func() Iterator

var stringType = reflect.TypeOf("")

func indirect(value reflect.Value) (reflect.Value, bool) {
	isNil := false
	for ; value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface; value = value.Elem() {
		if isNil = value.IsNil(); isNil {
			break
		}
	}
	return value, isNil
}

func neverConfirm(s string) (Object, bool) {
	return None, false
}

func loadFromStruct(value reflect.Value) ConvertString {
	return func(fieldName string) (Object, bool) {
		field := value.FieldByName(fieldName)
		if !field.IsValid() {
			return None, false
		}
		return valueToObject(field), true
	}
}

type structIterator struct {
	NoneType
	innerStruct reflect.Value
	innerType   reflect.Type
	current     int64
	end         int64
}

func (it *structIterator) Iter() Iterator {
	return it
}

func (it *structIterator) Next() (Object, bool) {
	var res Object = None
	ok := it.current < it.end
	if ok {
		current := int(it.current)
		name := it.innerType.Field(current).Name
		field := it.innerStruct.Field(current)
		res = NewList(String(name), valueToObject(field))
		it.current++
	}
	return res, ok
}

func (it *structIterator) Close() {
	it.innerStruct = reflect.Value{}
	it.innerType = nil
}

func iterFromStruct(value reflect.Value) ExtractIterator {
	return func() Iterator {
		return &structIterator{innerStruct: value, innerType: value.Type(), end: int64(value.NumField())}
	}
}

func loadFromMap(value reflect.Value) ConvertString {
	if !stringType.AssignableTo(value.Type().Key()) {
		return neverConfirm
	}
	return func(fieldName string) (Object, bool) {
		resValue := value.MapIndex(reflect.ValueOf(fieldName))
		if !resValue.IsValid() {
			return None, false
		}
		return valueToObject(resValue), true
	}
}

type mapIterWrapper struct {
	NoneType
	inner     *reflect.MapIter
	exhausted bool
}

func (w *mapIterWrapper) Iter() Iterator {
	return w
}

func (w *mapIterWrapper) Next() (Object, bool) {
	if !w.exhausted {
		if w.inner.Next() {
			return NewList(valueToObject(w.inner.Key()), valueToObject(w.inner.Value())), true
		} else {
			w.exhausted = true
		}
	}
	return None, false
}

func (w *mapIterWrapper) Close() {
	w.inner.Reset(reflect.Value{})
}

func iterFromMap(value reflect.Value) ExtractIterator {
	return func() Iterator {
		return &mapIterWrapper{inner: value.MapRange()}
	}
}

type dataWrapper struct {
	NoneType
	loadData ConvertString
	iterData ExtractIterator
}

func (w dataWrapper) Load(key Object) Object {
	return Load(w, key)
}

func (w dataWrapper) LoadStr(s string) (Object, bool) {
	return w.loadData(s)
}

func (w dataWrapper) Iter() Iterator {
	return w.iterData()
}

func valueToObject(value reflect.Value) Object {
	value, isNil := indirect(value)
	if !isNil {
		switch value.Kind() {
		case reflect.Bool:
			return Boolean(value.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Integer(value.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Integer(int64(value.Uint()))
		case reflect.Float32, reflect.Float64:
			return Float(value.Float())
		case reflect.Complex64, reflect.Complex128:
			return String(strconv.FormatComplex(value.Complex(), 'g', -1, 128))
		case reflect.String:
			return String(value.String())
		case reflect.Array, reflect.Slice:
			size := value.Len()
			l := &List{categories: map[string]NoneType{}, inner: make([]Object, 0, size)}
			for index := 0; index < size; index++ {
				l.Add(valueToObject(value.Index(index)))
			}
			return l
		case reflect.Struct:
			return dataWrapper{loadData: loadFromStruct(value), iterData: iterFromStruct(value)}
		case reflect.Map:
			return dataWrapper{loadData: loadFromMap(value), iterData: iterFromMap(value)}
		}
	}
	return None
}
