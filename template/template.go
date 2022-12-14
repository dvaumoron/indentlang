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
	local := types.NewLocalEnvironment(t.env)
	// TODO load data into local
	_, err := t.main.Apply(local, types.NewList()).WriteTo(w)
	return err
}

func Parse(s string) (*Template, error) {
	return ParseFrom("", s)
}

func ParseFrom(path, s string) (*Template, error) {
	baseEnv := builtins.Builtins
	if path != "" {
		baseEnv.Store(types.NewString("main"), builtins.NewImportDirective(path))
	}
	env := types.NewLocalEnvironment(baseEnv)

	parser.Parse(s).Eval(env)

	var tmpl *Template
	var err error
	main, exist := env.LoadConfirm(types.NewString("main"))
	if exist {
		casted, success := main.(types.Appliable)
		if success {
			tmpl = &Template{env: env, main: casted}
		} else {
			err = errors.New("the object main is not an Appliable")
		}
	} else {
		err = errors.New("the object main does not exist")
	}
	return tmpl, err
}
