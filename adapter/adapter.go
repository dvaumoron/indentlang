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
package adapter

import (
	"io/fs"
	"path/filepath"

	"github.com/dvaumoron/indentlang/builtins"
	"github.com/dvaumoron/indentlang/template"
)

func LoadTemplates(templatesPath string) (map[string]template.Template, error) {
	templatesPath, err := filepath.Abs(templatesPath)
	if err != nil {
		return nil, err
	}
	templatesPath = builtins.CheckPath(templatesPath)

	importDirective := builtins.MakeImportDirective(templatesPath)

	templates := map[string]template.Template{}
	inSize := len(templatesPath)
	err = filepath.WalkDir(templatesPath, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			name := path[inSize:]
			if end := len(name) - builtins.DefaultExtLen; name[end:] == builtins.DefaultExt {
				templates[name[:end]] = template.ParseWithImport(importDirective, name)
			}
		}
		return err
	})

	if err != nil {
		return nil, err
	}
	return templates, nil
}
