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
	"path/filepath"
	"strings"

	"github.com/dvaumoron/indentlang/builtins"
	"github.com/dvaumoron/indentlang/types"
)

type Template struct {
	env  types.Environment
	main types.Appliable
}

func (t *Template) Execute(w io.Writer, data any) error {
	// the LocalEnvironment layer is useful if Main is a macro
	dataEnv := types.NewLocalEnvironment(types.NewDataEnvironment(data, t.env))
	_, err := t.main.Apply(dataEnv, types.NewList()).WriteTo(w)
	return err
}

func ParsePath(path string) (*Template, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	splitIndex := strings.LastIndex(path, "/") + 1
	basePath, fileName := path[:splitIndex], path[splitIndex:]
	return ParseFrom(builtins.MakeImportDirective(basePath), fileName)
}

func ParseFrom(importDirective types.Appliable, filePath string) (*Template, error) {
	env := types.NewLocalEnvironment(builtins.Builtins)

	args := types.NewList()
	args.Add(types.NewString(filePath))
	importDirective.Apply(env, args)

	var tmpl *Template
	var err error
	main, ok := env.LoadStr(builtins.MainName)
	if ok {
		var mainAppliable types.Appliable
		mainAppliable, ok = main.(types.Appliable)
		if ok {
			tmpl = &Template{env: env, main: mainAppliable}
		} else {
			err = errors.New("the object Main is not an Appliable")
		}
	} else {
		err = errors.New("cannot load object Main")
	}
	return tmpl, err
}
