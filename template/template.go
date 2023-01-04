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
	env types.Environment
}

func (t *Template) Execute(w io.Writer, data any) error {
	var err error
	main, ok := t.env.LoadStr(builtins.MainName)
	if ok {
		var mainAppliable types.Appliable
		mainAppliable, ok = main.(types.Appliable)
		if ok {
			// each call must have its environment to avoid conflict in parallele execution
			local := types.NewLocalEnvironment(t.env)
			_, err = mainAppliable.ApplyWithData(data, local, types.NewList()).WriteTo(w)
		} else {
			err = errors.New("the object Main is not an Appliable")
		}
	} else {
		err = errors.New("cannot load object Main")
	}
	return err
}

func ParsePath(path string) (*Template, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	splitIndex := strings.LastIndex(path, "/") + 1
	basePath, fileName := path[:splitIndex], path[splitIndex:]
	return ParseWithImport(builtins.MakeImportDirective(basePath), fileName), nil
}

func ParseWithImport(importDirective types.Appliable, filePath string) *Template {
	env := types.NewLocalEnvironment(builtins.Builtins)
	env.StoreStr(builtins.ImportName, importDirective)

	importDirective.Apply(env, types.NewList(types.String(filePath)))

	return &Template{env: env}
}
