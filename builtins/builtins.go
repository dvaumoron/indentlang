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

const MainName = "Main"

var Builtins = initBuitins()

func initBuitins() types.BaseEnvironment {
	elementHtml := CreateHtmlElement("html")

	base := types.MakeBaseEnvironment()
	// special case in order to create Main,
	// which will be called by the Execute method of the Template struct
	base.StoreStr("html", types.MakeNativeAppliable(func(env types.Environment, args types.Iterable) types.Object {
		env.StoreStr(MainName, types.MakeNativeAppliable(func(callEnv types.Environment, emptyArgs types.Iterable) types.Object {
			return elementHtml.Apply(callEnv, args)
		}))
		return types.None
	}))
	// all other not deprecated html element
	addHtmlElement(base, "a")
	addHtmlElement(base, "abbr")
	addHtmlElement(base, "address")
	addHtmlElement(base, "area")
	addHtmlElement(base, "article")
	addHtmlElement(base, "aside")
	addHtmlElement(base, "audio")
	addHtmlElement(base, "b")
	addHtmlElement(base, "base")
	addHtmlElement(base, "bdi")
	addHtmlElement(base, "bdo")
	addHtmlElement(base, "blockquote")
	addHtmlElement(base, "body")
	addHtmlElement(base, "br")
	addHtmlElement(base, "button")
	addHtmlElement(base, "canvas")
	addHtmlElement(base, "caption")
	addHtmlElement(base, "cite")
	addHtmlElement(base, "code")
	addHtmlElement(base, "col")
	addHtmlElement(base, "colgroup")
	addHtmlElement(base, "data")
	addHtmlElement(base, "datalist")
	addHtmlElement(base, "dd")
	addHtmlElement(base, "del")
	addHtmlElement(base, "details")
	addHtmlElement(base, "dfn")
	addHtmlElement(base, "dialog")
	addHtmlElement(base, "div")
	addHtmlElement(base, "dl")
	addHtmlElement(base, "dt")
	addHtmlElement(base, "em")
	addHtmlElement(base, "embed")
	addHtmlElement(base, "fieldset")
	addHtmlElement(base, "figcaption")
	addHtmlElement(base, "figure")
	addHtmlElement(base, "footer")
	addHtmlElement(base, "form")
	addHtmlElement(base, "h1")
	addHtmlElement(base, "head")
	addHtmlElement(base, "header")
	addHtmlElement(base, "hgroup")
	addHtmlElement(base, "hr")
	addHtmlElement(base, "i")
	addHtmlElement(base, "iframe")
	addHtmlElement(base, "img")
	addHtmlElement(base, "input")
	addHtmlElement(base, "ins")
	addHtmlElement(base, "kbd")
	addHtmlElement(base, "label")
	addHtmlElement(base, "legend")
	addHtmlElement(base, "li")
	addHtmlElement(base, "link")
	addHtmlElement(base, "main")
	addHtmlElement(base, "map")
	addHtmlElement(base, "mark")
	addHtmlElement(base, "menu")
	addHtmlElement(base, "meta")
	addHtmlElement(base, "meter")
	addHtmlElement(base, "nav")
	addHtmlElement(base, "noscript")
	addHtmlElement(base, "object")
	addHtmlElement(base, "ol")
	addHtmlElement(base, "optgroup")
	addHtmlElement(base, "option")
	addHtmlElement(base, "output")
	addHtmlElement(base, "p")
	addHtmlElement(base, "picture")
	addHtmlElement(base, "pre")
	addHtmlElement(base, "progress")
	addHtmlElement(base, "q")
	addHtmlElement(base, "rp")
	addHtmlElement(base, "rt")
	addHtmlElement(base, "ruby")
	addHtmlElement(base, "s")
	addHtmlElement(base, "samp")
	addHtmlElement(base, "script")
	addHtmlElement(base, "section")
	addHtmlElement(base, "select")
	addHtmlElement(base, "slot")
	addHtmlElement(base, "small")
	addHtmlElement(base, "source")
	addHtmlElement(base, "span")
	addHtmlElement(base, "strong")
	addHtmlElement(base, "style")
	addHtmlElement(base, "sub")
	addHtmlElement(base, "summary")
	addHtmlElement(base, "sup")
	addHtmlElement(base, "table")
	addHtmlElement(base, "tbody")
	addHtmlElement(base, "td")
	addHtmlElement(base, "template")
	addHtmlElement(base, "textarea")
	addHtmlElement(base, "tfoot")
	addHtmlElement(base, "th")
	addHtmlElement(base, "thead")
	addHtmlElement(base, "time")
	addHtmlElement(base, "title")
	addHtmlElement(base, "tr")
	addHtmlElement(base, "track")
	addHtmlElement(base, "u")
	addHtmlElement(base, "ul")
	addHtmlElement(base, "var")
	addHtmlElement(base, "video")
	addHtmlElement(base, "wbr")

	// true langage features
	// *Form indicate a special form
	// *Func indicate a normal function
	base.StoreStr("If", types.MakeNativeAppliable(ifForm))
	base.StoreStr("For", types.MakeNativeAppliable(forForm))
	base.StoreStr("While", types.MakeNativeAppliable(whileForm))
	base.StoreStr(parser.SetName, types.MakeNativeAppliable(setForm))
	base.StoreStr(".", types.MakeNativeAppliable(getForm))
	base.StoreStr("[]", types.MakeNativeAppliable(loadForm))
	base.StoreStr("[]=", types.MakeNativeAppliable(storeFunc))
	base.StoreStr("Func", types.MakeNativeAppliable(funcForm))
	base.StoreStr("Lambda", types.MakeNativeAppliable(lambdaForm))
	base.StoreStr("Call", types.MakeNativeAppliable(callForm))
	base.StoreStr("Macro", types.MakeNativeAppliable(macroForm))
	// TODO init stuff
	// Quote, Unquote, type conversion
	// List, Dict, Range, Enumerate, Add, Size, Del, Iter, Next, AddAttribute, HasAttribute
	// And, Or, Not, ==, !=, >, >=, <, <=, +, -, *, /, //, %

	// give parser package a protected copy to use in user custom rules
	parser.BuiltinsCopy = types.NewLocalEnvironment(base)
	return base
}

func addHtmlElement(base types.BaseEnvironment, name string) {
	base.StoreStr(name, CreateHtmlElement(name))
}

var openElement = types.NewString("<")
var openCloseElement = types.NewString("</")
var closeElement = types.NewString(">")
var closeOpenElement = types.NewString("/>")
var space = types.NewString(" ")
var equalQuote = types.NewString("=\"")
var quote = types.NewString("\"")

func CreateHtmlElement(name string) types.NativeAppliable {
	wrappedName := types.NewString(name)
	return types.MakeNativeAppliable(func(env types.Environment, args types.Iterable) types.Object {
		local := types.NewLocalEnvironment(env)
		attrs := types.NewList()
		childs := types.NewList()
		types.ForEach(args, func(value types.Object) bool {
			value = value.Eval(local)
			if value.HasCategory(parser.AttributeName) {
				attrs.Add(value)
			} else {
				childs.Add(value)
			}
			return true
		})
		res := types.NewList()
		res.Add(openElement)
		res.Add(wrappedName)
		types.ForEach(attrs, func(value types.Object) bool {
			attr, success := value.(types.Iterable)
			if success {
				itAttr := attr.Iter()
				attrName, exist := itAttr.Next()
				if exist {
					res.Add(space)
					res.Add(attrName)
					attrValue, exist := itAttr.Next()
					if exist {
						res.Add(equalQuote)
						res.Add(attrValue)
						res.Add(quote)
					}
				}
			}
			return true
		})
		if childs.SizeInt() == 0 {
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
