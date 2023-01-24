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
const UnquoteName = "Unquote"

var CustomRules = types.NewList()

var wordParsers []types.ConvertString

// an empty environment to execute custom rules
var BuiltinsCopy types.Environment = types.MakeBaseEnvironment()

// needed to prevent a cycle in the initialisation
func init() {
	wordParsers = []types.ConvertString{
		parseTrue, parseFalse, parseNone, parseAttribute, parseUnquote,
		parseList, parseString, parseString2, parseInt, parseFloat,
	}
}

func HandleClassicWord(word string, nodeList *types.List) {
	if nativeRules(word, nodeList) {
		args := types.NewList(types.String(word))
		res := true
		types.ForEach(CustomRules, func(object types.Object) bool {
			rule, ok := object.(types.Appliable)
			if ok {
				// The Apply must return None if it fails.
				node := rule.Apply(BuiltinsCopy, args)
				_, res = node.(types.NoneType)
				if !res {
					nodeList.Add(node)
				}
			}
			return res
		})
		if res {
			nodeList.Add(types.Identifier(word))
		}
	}
}

// a true is returned when no rule match
func nativeRules(word string, nodeList *types.List) bool {
	for _, parser := range wordParsers {
		node, ok := parser(word)
		if ok {
			nodeList.Add(node)
			return false
		}
	}
	return true
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
	lastIndex := len(word) - 1
	if word[0] != '"' || word[lastIndex] != '"' {
		return nil, false
	}
	escape := false
	extracted := make([]rune, 0, lastIndex)
	for _, char := range word[1:lastIndex] {
		if escape {
			escape = false
			if char == '\'' {
				extracted = append(extracted, '\'')
			} else {
				extracted = append(extracted, '\\', char)
			}
		} else {
			switch char {
			case '"':
				return nil, false
			case '\\':
				escape = true
			default:
				extracted = append(extracted, char)
			}
		}
	}
	return types.String(extracted), true
}

func parseString2(word string) (types.Object, bool) {
	lastIndex := len(word) - 1
	if word[0] != '\'' || word[lastIndex] != '\'' {
		return nil, false
	}
	escape := false
	extracted := make([]rune, 0, lastIndex)
	for _, char := range word[1:lastIndex] {
		if escape {
			escape = false
			if char == '\'' {
				extracted = append(extracted, '\'')
			} else {
				extracted = append(extracted, '\\', char)
			}
		} else {
			switch char {
			case '"':
				extracted = append(extracted, '\\', '"')
			case '\'':
				return nil, false
			case '\\':
				escape = true
			default:
				extracted = append(extracted, char)
			}
		}
	}
	return types.String(extracted), true
}

func parseAttribute(word string) (types.Object, bool) {
	if word[0] != '@' {
		return nil, false
	}
	elems := strings.SplitN(word[1:], "=", 2)
	attr := types.NewList(types.String(elems[0]))
	attr.AddCategory(AttributeName)
	if len(elems) > 1 {
		HandleClassicWord(elems[1], attr)
	}
	return attr, true
}

// manage melting with string literal
func parseList(word string) (types.Object, bool) {
	if word == SetName {
		return nil, false
	}
	chars := make(chan rune)
	go sendChar(chars, word)
	index := 0
	var indexes []int
	for char := range chars {
		switch char {
		case '"', '\'':
			delim := char
		InnerCharLoop:
			for char := range chars {
				index++
				switch char {
				case delim:
					// no need of unended string detection,
					// this have already been tested in the word splitting part
					break InnerCharLoop
				case '\\':
					<-chars
					index++
				}
			}
		case ':':
			indexes = append(indexes, index)
		}
		index++
	}
	if len(indexes) == 0 {
		return nil, false
	}
	nodeList := types.NewList(ListId)
	startIndex := 0
	for _, splitIndex := range indexes {
		handleSubWord(word[startIndex:splitIndex], nodeList)
		startIndex = splitIndex + 1
	}
	handleSubWord(word[startIndex:], nodeList)
	return nodeList, true
}

func handleSubWord(word string, nodeList *types.List) {
	if word == "" {
		nodeList.Add(types.None)
	} else {
		HandleClassicWord(word, nodeList)
	}
}

func parseInt(word string) (types.Object, bool) {
	i, err := strconv.ParseInt(word, 10, 64)
	return types.Integer(i), err == nil
}

func parseFloat(word string) (types.Object, bool) {
	f, err := strconv.ParseFloat(word, 64)
	return types.Float(f), err == nil
}

func parseUnquote(word string) (types.Object, bool) {
	if word[0] != ',' {
		return nil, false
	}
	nodeList := types.NewList(types.Identifier(UnquoteName))
	handleSubWord(word[1:], nodeList)
	return nodeList, true
}
