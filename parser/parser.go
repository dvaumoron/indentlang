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
package parser

import (
	"errors"
	"io"

	"github.com/dvaumoron/indentlang/types"
)

type Template struct {
	main types.Function
}

func (*Template) Execute(w io.Writer, data any) error {
	// TODO
	return nil
}

func Parse(s string) (*Template, error) {
	var err error
	var env types.Environment
	// TODO

	var tmpl *Template

	mainObject := env.Get(types.NewString("main"))
	mainFunction, success := mainObject.(types.Function)
	if success {
		tmpl = &Template{main: mainFunction}
	} else {
		err = errors.New("The object main is not a Function.")
	}
	return tmpl, err
}
