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
	"reflect"
)

var stringType = reflect.TypeOf("")

func indirect(value reflect.Value) (reflect.Value, bool) {
	for ; value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface; value = value.Elem() {
		if value.IsNil() {
			return value, true
		}
	}
	return value, false
}

func neverConfirm(s string) (Object, bool) {
	return None, false
}

func loadFromStruct(value reflect.Value) func(string) (Object, bool) {
	vType := value.Type()
	return func(fieldName string) (Object, bool) {
		var res Object = None
		fType, ok := vType.FieldByName(fieldName)
		if ok && fType.IsExported() {
			field, err := value.FieldByIndexErr(fType.Index)
			if err == nil {
				res = valueToObject(field)
			} else {
				ok = false
			}
		}
		return res, ok
	}
}

func loadFromMap(value reflect.Value) func(string) (Object, bool) {
	return func(fieldName string) (Object, bool) {
		var res Object = None
		resValue := value.MapIndex(reflect.ValueOf(fieldName))
		ok := resValue.IsValid()
		if ok {
			res = valueToObject(resValue)
		}
		return res, ok
	}
}

type LoadWrapper struct {
	NoneType
	loadConfirm func(string) (Object, bool)
}

func (w LoadWrapper) Load(key Object) Object {
	var res Object = None
	str, ok := key.(*String)
	if ok {
		res, _ = w.loadConfirm(str.Inner)
	}
	return res
}

func (w LoadWrapper) LoadStr(s string) Object {
	res, _ := w.loadConfirm(s)
	return res
}

func valueToObject(value reflect.Value) Object {
	var res Object = None
	value, isNil := indirect(value)
	if !isNil {
		switch value.Kind() {
		case reflect.Bool:
			if value.Bool() {
				res = True
			} else {
				res = False
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			res = NewInteger(value.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			res = NewInteger(int64(value.Uint()))
		case reflect.Float32, reflect.Float64:
			res = NewFloat(value.Float())
		case reflect.Complex64, reflect.Complex128:
			res = NewString(fmt.Sprint(value.Complex()))
		case reflect.String:
			res = NewString(value.String())
		case reflect.Array, reflect.Slice:
			size := value.Len()
			l := &List{categories: makeCategories(), inner: make([]Object, 0, size)}
			for index := 0; index < size; index++ {
				l.Add(valueToObject(value.Index(index)))
			}
			res = l
		case reflect.Struct:
			res = LoadWrapper{loadConfirm: loadFromStruct(value)}
		case reflect.Map:
			if stringType.AssignableTo(value.Type().Key()) {
				res = LoadWrapper{loadConfirm: loadFromMap(value)}
			}
		}
	}
	return res
}
