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
package main

import (
	"fmt"
	"os"

	"github.com/dvaumoron/indentlang/template"
	"gopkg.in/yaml.v3"
)

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("Usage : indentlang file.il data.yaml outputFile")
		return
	}

	tmplPath := args[1]
	dataPath := args[2]
	outPath := args[3]

	tmpl, err := template.ParsePath(tmplPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	dataBody, err := os.ReadFile(dataPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	tmplArgs := map[string]any{}
	err = yaml.Unmarshal(dataBody, tmplArgs)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create(outPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	err = tmpl.Execute(file, tmplArgs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(outPath, "generated")
}
