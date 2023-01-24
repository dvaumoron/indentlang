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
package builtins

import (
	"github.com/dvaumoron/indentlang/parser"
	"github.com/dvaumoron/indentlang/types"
)

const openElement types.String = "<"
const openCloseElement types.String = "</"
const closeElement types.String = ">"
const closeOpenElement types.String = "/>"
const space types.String = " "
const equalQuote types.String = "=\""
const quote types.String = "\""

func addHtmlElement(base types.BaseEnvironment, name string) {
	base.StoreStr(name, CreateHtmlElement(name))
}

func CreateHtmlElement(name string) types.NativeAppliable {
	wrappedName := types.String(name)
	return types.MakeNativeAppliable(func(env types.Environment, itArgs types.Iterator) types.Object {
		attrs := types.NewList()
		childs := types.NewList()
		types.ForEach(itArgs, func(arg types.Object) bool {
			switch casted := arg.Eval(env).(type) {
			case types.NoneType:
				// ignore None
			case *types.List:
				if casted.HasCategory(parser.AttributeName) {
					attrs.Add(casted)
				} else {
					childs.Add(casted)
				}
			default:
				childs.Add(casted)
			}
			return true
		})
		res := types.NewList(openElement, wrappedName)
		types.ForEach(attrs, func(value types.Object) bool {
			attr, ok := value.(types.Iterable)
			if !ok {
				return true
			}

			itAttr := attr.Iter()
			defer itAttr.Close()
			attrName, ok := itAttr.Next()
			if !ok {
				return true
			}

			res.Add(space)
			res.Add(attrName)
			attrValue, ok := itAttr.Next()
			if ok {
				res.Add(equalQuote)
				res.Add(attrValue)
				res.Add(quote)
			}
			return true
		})
		if childs.Size() == 0 {
			res.Add(closeOpenElement)
		} else {
			res.Add(closeElement)
			types.ForEach(childs, func(value types.Object) bool {
				res.Add(space)
				res.Add(value)
				return true
			})
			res.Add(openCloseElement)
			res.Add(wrappedName)
			res.Add(closeElement)
		}
		return res
	})
}
