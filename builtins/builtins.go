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
	elementHtml := createXmlTag("html")

	base := types.MakeBaseEnvironment()
	// special case in order to create Main,
	// which will be called by the Execute method of the Template struct
	base.StoreStr("html", types.MakeNativeAppliable(func(env types.Environment, itArgs types.Iterator) types.Object {
		// avoid loss in multiple call case
		savedArgs := types.NewList().AddAll(itArgs)
		env.StoreStr(MainName, types.MakeNativeAppliable(func(callEnv types.Environment, emptyArgs types.Iterator) types.Object {
			return elementHtml.Apply(callEnv, savedArgs)
		}))
		return types.None
	}))
	// all other not deprecated html element
	addXmlTag(base, "a")
	addXmlTag(base, "abbr")
	addXmlTag(base, "address")
	addXmlTag(base, "area")
	addXmlTag(base, "article")
	addXmlTag(base, "aside")
	addXmlTag(base, "audio")
	addXmlTag(base, "b")
	addXmlTag(base, "base")
	addXmlTag(base, "bdi")
	addXmlTag(base, "bdo")
	addXmlTag(base, "blockquote")
	addXmlTag(base, "body")
	addXmlTag(base, "br")
	addXmlTag(base, "button")
	addXmlTag(base, "canvas")
	addXmlTag(base, "caption")
	addXmlTag(base, "cite")
	addXmlTag(base, "code")
	addXmlTag(base, "col")
	addXmlTag(base, "colgroup")
	addXmlTag(base, "data")
	addXmlTag(base, "datalist")
	addXmlTag(base, "dd")
	addXmlTag(base, "del")
	addXmlTag(base, "details")
	addXmlTag(base, "dfn")
	addXmlTag(base, "dialog")
	addXmlTag(base, "div")
	addXmlTag(base, "dl")
	addXmlTag(base, "dt")
	addXmlTag(base, "em")
	addXmlTag(base, "embed")
	addXmlTag(base, "fieldset")
	addXmlTag(base, "figcaption")
	addXmlTag(base, "figure")
	addXmlTag(base, "footer")
	addXmlTag(base, "form")
	addXmlTag(base, "h1")
	addXmlTag(base, "h2")
	addXmlTag(base, "h3")
	addXmlTag(base, "h4")
	addXmlTag(base, "h5")
	addXmlTag(base, "h6")
	addXmlTag(base, "head")
	addXmlTag(base, "header")
	addXmlTag(base, "hgroup")
	addXmlTag(base, "hr")
	addXmlTag(base, "i")
	addXmlTag(base, "iframe")
	addXmlTag(base, "img")
	addXmlTag(base, "input")
	addXmlTag(base, "ins")
	addXmlTag(base, "kbd")
	addXmlTag(base, "label")
	addXmlTag(base, "legend")
	addXmlTag(base, "li")
	addXmlTag(base, "link")
	addXmlTag(base, "main")
	addXmlTag(base, "map")
	addXmlTag(base, "mark")
	addXmlTag(base, "menu")
	addXmlTag(base, "meta")
	addXmlTag(base, "meter")
	addXmlTag(base, "nav")
	addXmlTag(base, "noscript")
	addXmlTag(base, "object")
	addXmlTag(base, "ol")
	addXmlTag(base, "optgroup")
	addXmlTag(base, "option")
	addXmlTag(base, "output")
	addXmlTag(base, "p")
	addXmlTag(base, "picture")
	addXmlTag(base, "pre")
	addXmlTag(base, "progress")
	addXmlTag(base, "q")
	addXmlTag(base, "rp")
	addXmlTag(base, "rt")
	addXmlTag(base, "ruby")
	addXmlTag(base, "s")
	addXmlTag(base, "samp")
	addXmlTag(base, "script")
	addXmlTag(base, "section")
	addXmlTag(base, "select")
	addXmlTag(base, "slot")
	addXmlTag(base, "small")
	addXmlTag(base, "source")
	addXmlTag(base, "span")
	addXmlTag(base, "strong")
	addXmlTag(base, "style")
	addXmlTag(base, "sub")
	addXmlTag(base, "summary")
	addXmlTag(base, "sup")
	addXmlTag(base, "table")
	addXmlTag(base, "tbody")
	addXmlTag(base, "td")
	addXmlTag(base, "template")
	addXmlTag(base, "textarea")
	addXmlTag(base, "tfoot")
	addXmlTag(base, "th")
	addXmlTag(base, "thead")
	addXmlTag(base, "time")
	addXmlTag(base, "title")
	addXmlTag(base, "tr")
	addXmlTag(base, "track")
	addXmlTag(base, "u")
	addXmlTag(base, "ul")
	addXmlTag(base, "var")
	addXmlTag(base, "video")
	addXmlTag(base, "wbr")

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
	// allowing other XMLs beyond HTML
	base.StoreStr("XmlTag", types.MakeNativeAppliable(xmlTagFunc))

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
	base.StoreStr("Close", types.MakeNativeAppliable(closeFunc))
	base.StoreStr("Size", types.MakeNativeAppliable(sizeFunc))
	base.StoreStr("Add", types.MakeNativeAppliable(addFunc))
	base.StoreStr("AddAll", types.MakeNativeAppliable(addAllFunc))

	// some functions to do math
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
	base.StoreStr("GetEnv", types.MakeNativeAppliable(getEnvFunc))

	// escape functions
	base.StoreStr("EscapeHtml", types.MakeNativeAppliable(escapeHtmlFunc))
	base.StoreStr("EscapeQuery", types.MakeNativeAppliable(escapeQueryFunc))
	base.StoreStr("EscapePath", types.MakeNativeAppliable(escapePathFunc))

	// TODO init stuff
	// lack of utilities (for iterator, function, ...)

	// give parser package a protected copy to use in user custom rules
	parser.BuiltinsCopy = types.MakeLocalEnvironment(base)
	return base
}
