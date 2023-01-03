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
	base.StoreStr("html", types.MakeNativeAppliable(func(env types.Environment, itArgs types.Iterator) types.Object {
		env.StoreStr(MainName, types.MakeNativeAppliable(func(callEnv types.Environment, emptyArgs types.Iterator) types.Object {
			return elementHtml.Apply(callEnv, itArgs)
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

	// start of the "true" language features
	// *Form indicate a special form
	// *Func indicate a normal function
	base.StoreStr("If", types.MakeNativeAppliable(ifForm))
	base.StoreStr("For", types.MakeNativeAppliable(forForm))
	base.StoreStr("While", types.MakeNativeAppliable(whileForm))
	base.StoreStr(parser.SetName, types.MakeNativeAppliable(setForm))
	base.StoreStr(".", types.MakeNativeAppliable(getForm))
	base.StoreStr("[]", types.MakeNativeAppliable(loadFunc))
	base.StoreStr("[]=", types.MakeNativeAppliable(storeFunc))

	// functions and macros management
	base.StoreStr("Func", types.MakeNativeAppliable(funcForm))
	base.StoreStr("Lambda", types.MakeNativeAppliable(lambdaForm))
	base.StoreStr("Call", types.MakeNativeAppliable(callFunc))
	base.StoreStr("Macro", types.MakeNativeAppliable(macroForm))

	// basic type creation and conversion
	base.StoreStr("Identifier", types.MakeNativeAppliable(identifierConvFunc))
	base.StoreStr("Boolean", types.MakeNativeAppliable(boolConvFunc))
	base.StoreStr("Integer", types.MakeNativeAppliable(intConvFunc))
	base.StoreStr("Float", types.MakeNativeAppliable(floatConvFunc))
	base.StoreStr("String", types.MakeNativeAppliable(stringConvFunc))
	base.StoreStr(string(parser.ListId), types.MakeNativeAppliable(listFunc))
	base.StoreStr("Dict", types.MakeNativeAppliable(dictFunc))

	// logic
	base.StoreStr("Not", types.MakeNativeAppliable(notFunc))
	base.StoreStr("And", types.MakeNativeAppliable(andFunc))
	base.StoreStr("Or", types.MakeNativeAppliable(orFunc))
	base.StoreStr("==", types.MakeNativeAppliable(equalsFunc))
	base.StoreStr("!=", types.MakeNativeAppliable(notEqualsFunc))
	base.StoreStr(">", types.MakeNativeAppliable(greaterThanFunc))
	base.StoreStr(">=", types.MakeNativeAppliable(greaterEqualFunc))
	base.StoreStr("<", types.MakeNativeAppliable(lessThanFunc))
	base.StoreStr("<=", types.MakeNativeAppliable(lessEqualFunc))

	// advanced looping
	base.StoreStr("Range", types.MakeNativeAppliable(rangeFunc))
	base.StoreStr("Enumerate", types.MakeNativeAppliable(enumerateFunc))
	base.StoreStr("Iter", types.MakeNativeAppliable(iterFunc))
	base.StoreStr("Next", types.MakeNativeAppliable(nextFunc))
	base.StoreStr("Size", types.MakeNativeAppliable(sizeFunc))
	base.StoreStr("Add", types.MakeNativeAppliable(addFunc))
	base.StoreStr("AddAll", types.MakeNativeAppliable(addAllFunc))

	// some function to do math
	base.StoreStr(sumName, types.MakeNativeAppliable(sumFunc))
	base.StoreStr(minusName, types.MakeNativeAppliable(minusFunc))
	base.StoreStr(productName, types.MakeNativeAppliable(productFunc))
	base.StoreStr(divideName, types.MakeNativeAppliable(divideFunc))
	base.StoreStr(floorDivideName, types.MakeNativeAppliable(floorDivideFunc))
	base.StoreStr(remainderName, types.MakeNativeAppliable(remainderFunc))
	base.StoreStr("+=", types.MakeNativeAppliable(sumSetForm))
	base.StoreStr("-=", types.MakeNativeAppliable(minusSetForm))
	base.StoreStr("*=", types.MakeNativeAppliable(productSetForm))
	base.StoreStr("/=", types.MakeNativeAppliable(divideSetForm))
	base.StoreStr("//=", types.MakeNativeAppliable(floorDivideSetForm))
	base.StoreStr("%=", types.MakeNativeAppliable(remainderSetForm))

	// advanced programming
	base.StoreStr("Quote", types.MakeNativeAppliable(quoteForm))
	base.StoreStr(parser.UnquoteName, types.MakeNativeAppliable(unquoteFunc)) // not very useful
	base.StoreStr("Eval", types.MakeNativeAppliable(evalForm))
	base.StoreStr("Del", types.MakeNativeAppliable(delForm))
	base.StoreStr("AddCategory", types.MakeNativeAppliable(addCategoryFunc))
	base.StoreStr("HasCategory", types.MakeNativeAppliable(hasCategoryFunc))
	base.StoreStr("AddCustomRule", types.MakeNativeAppliable(addCustomRuleFunc))
	base.StoreStr("ParseWord", types.MakeNativeAppliable(parseWordFunc))

	// TODO init stuff
	// lack of utilities (for string, iterator, function, ...)

	// give parser package a protected copy to use in user custom rules
	parser.BuiltinsCopy = types.NewLocalEnvironment(base)
	return base
}
