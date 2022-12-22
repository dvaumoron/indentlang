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
package template

import (
	"errors"
	"io"

	"github.com/dvaumoron/indentlang/builtins"
	"github.com/dvaumoron/indentlang/parser"
	"github.com/dvaumoron/indentlang/types"
)

type Template struct {
	env  types.Environment
	main types.Appliable
}

func (t *Template) Execute(w io.Writer, data any) error {
	local := types.NewDataEnvironment(data, types.NewLocalEnvironment(t.env))
	_, err := t.main.Apply(local, types.NewList()).WriteTo(w)
	return err
}

func Parse(str string) (*Template, error) {
	return ParseFrom("", str)
}

func ParseFrom(path, str string) (*Template, error) {
	env := types.NewLocalEnvironment(builtins.Builtins)
	if path != "" {
		env.StoreStr("Import", builtins.MakeImportDirective(path))
	}

	var tmpl *Template
	node, err := parser.Parse(str)
	if err == nil {
		node.Eval(env)

		main, success := env.LoadStr("Main").(types.Appliable)
		if success {
			tmpl = &Template{env: env, main: main}
		} else {
			err = errors.New("the object Main is not an Appliable")
		}
	}
	return tmpl, err
}
