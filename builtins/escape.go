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
	"html"
	"net/url"

	"github.com/dvaumoron/indentlang/types"
)

func escapeHtmlFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return escapingFunc(env, itArgs, html.EscapeString)
}

func escapeQueryFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return escapingFunc(env, itArgs, url.QueryEscape)
}

func escapePathFunc(env types.Environment, itArgs types.Iterator) types.Object {
	return escapingFunc(env, itArgs, url.PathEscape)
}

func escapingFunc(env types.Environment, itArgs types.Iterator, escapeFunction func(string) string) types.Object {
	arg, _ := itArgs.Next()
	str, _ := arg.Eval(env).(types.String)
	if str == "" {
		return types.None
	}
	return types.String(escapeFunction(string(str)))
}
