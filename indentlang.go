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
	"strconv"

	"github.com/dvaumoron/indentlang/template"
)

func main() {
	args := os.Args

	tmplPath := args[1]
	outPath := args[2]

	tmplBody, err := os.ReadFile(tmplPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	tmpl, err := template.Parse(string(tmplBody))
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

	tmplArgs := map[string]any{}
	values := args[3:]
	tmplArgs["values"] = values
	for i, value := range values {
		tmplArgs["value"+strconv.Itoa(i+1)] = value
	}

	err = tmpl.Execute(file, tmplArgs)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(outPath, "generated")
}
