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
	"strconv"
	"strings"

	"github.com/dvaumoron/indentlang/types"
)

const SetName = ":="

var CustomRules = types.NewList()

var wordParsers []types.ConvertString

// an empty environment to execute custom rules
var BuiltinsCopy types.Environment = types.MakeBaseEnvironment()

// needed to prevent a cycle in the initialisation
func init() {
	wordParsers = []types.ConvertString{
		parseTrue, parseFalse, parseNone, parseString, parseString2,
		parseAttribute, parseList, parseInt, parseFloat,
	}
}

func handleWord(words <-chan string, listStack *stack[*types.List], done chan<- types.NoneType) {
	for word := range words {
		if word == "(" {
			manageOpen(listStack)
		} else if word == ")" {
			listStack.pop()
		} else {
			HandleClassicWord(word, listStack.peek())
		}
	}
	done <- types.None
}

func HandleClassicWord(word string, list *types.List) {
	if !nativeRules(word, list) {
		args := types.NewList(types.String(word))
		res := true
		types.ForEach(CustomRules, func(object types.Object) bool {
			rule, ok := object.(types.Appliable)
			if ok {
				// The Apply must return None if it fails.
				node := rule.Apply(BuiltinsCopy, args)
				_, res = node.(types.NoneType)
				if !res {
					list.Add(node)
				}
			}
			return res
		})
		if res {
			list.Add(types.Identifier(word))
		}
	}
}

func nativeRules(word string, list *types.List) bool {
	ok := false
	for _, parser := range wordParsers {
		var node types.Object
		node, ok = parser(word)
		if ok {
			list.Add(node)
			break
		}
	}
	return ok
}

func parseTrue(word string) (types.Object, bool) {
	return types.Boolean(true), word == "true"
}

func parseFalse(word string) (types.Object, bool) {
	return types.Boolean(false), word == "false"
}

func parseNone(word string) (types.Object, bool) {
	return types.None, word == "None"
}

func parseString(word string) (types.Object, bool) {
	var res types.Object
	ok := word[0] == '"'
	if ok {
		extracted := strings.ReplaceAll(word[1:len(word)-1], "\\'", "'")
		res = types.String(extracted)
	}
	return res, ok
}

func parseString2(word string) (types.Object, bool) {
	var res types.Object
	ok := word[0] == '\''
	if ok {
		skipApos := false
		size := len(word)
		extracted := make([]rune, 0, size)
		for _, char := range word[1 : size-1] {
			if skipApos {
				skipApos = false
				if char == '\'' {
					extracted = append(extracted, char)
				} else {
					extracted = append(extracted, '\\', char)
				}
			} else if char == '"' {
				extracted = append(extracted, '\\', char)
			} else if char == '\\' {
				skipApos = true
			} else {
				extracted = append(extracted, char)
			}
		}
		res = types.String(extracted)
	}
	return res, ok
}

func parseAttribute(word string) (types.Object, bool) {
	var res types.Object
	ok := word[0] == '@'
	if ok {
		elems := strings.Split(word[1:], "=")
		attr := types.NewList(types.String(elems[0]))
		attr.AddCategory(AttributeName)
		if len(elems) > 1 {
			HandleClassicWord(elems[1], attr)
		}
		res = attr
	}
	return res, ok
}

// TODO improve melting with string literal
func parseList(word string) (types.Object, bool) {
	var res types.Object = types.None
	ok := word != SetName && strings.Contains(word, ":")
	if ok {
		nodeList := types.NewList(ListId)
		for _, elem := range strings.Split(word, ":") {
			if elem == "" {
				nodeList.Add(types.None)
			} else {
				HandleClassicWord(elem, nodeList)
			}
		}
		res = nodeList
	}
	return res, ok
}

func parseInt(word string) (types.Object, bool) {
	i, err := strconv.ParseInt(word, 10, 64)
	return types.Integer(i), err == nil
}

func parseFloat(word string) (types.Object, bool) {
	f, err := strconv.ParseFloat(word, 64)
	return types.Float(f), err == nil
}
