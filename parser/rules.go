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
		parseTrue, parseFalse, parseNone, parseAttribute, parseList,
		parseString, parseString2, parseInt, parseFloat,
	}
}

func handleWord(words <-chan string, listStack *stack[*types.List], done chan<- types.NoneType) {
	for word := range words {
		switch word {
		case "(":
			manageOpen(listStack)
		case ")":
			listStack.pop()
		default:
			HandleClassicWord(word, listStack.peek())
		}
	}
	done <- types.None
}

func HandleClassicWord(word string, nodeList *types.List) {
	if !nativeRules(word, nodeList) {
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

func nativeRules(word string, nodeList *types.List) bool {
	ok := false
	for _, parser := range wordParsers {
		var node types.Object
		node, ok = parser(word)
		if ok {
			nodeList.Add(node)
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
	ok := word[0] == '"' && word[len(word)-1] == '"'
	if ok {
		extracted := strings.ReplaceAll(word[1:len(word)-1], "\\'", "'")
		res = types.String(extracted)
	}
	return res, ok
}

func parseString2(word string) (types.Object, bool) {
	var res types.Object
	ok := word[0] == '\'' && word[len(word)-1] == '\''
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
		elems := strings.SplitN(word[1:], "=", 2)
		attr := types.NewList(types.String(elems[0]))
		attr.AddCategory(AttributeName)
		if len(elems) > 1 {
			HandleClassicWord(elems[1], attr)
		}
		res = attr
	}
	return res, ok
}

// manage melting with string literal
func parseList(word string) (types.Object, bool) {
	var res types.Object = types.None
	ok := word != SetName
	if ok {
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
		ok = len(indexes) != 0
		if ok {
			nodeList := types.NewList(ListId)
			startIndex := 0
			for _, splitIndex := range indexes {
				handleSubWord(word[startIndex:splitIndex], nodeList)
				startIndex = splitIndex + 1
			}
			handleSubWord(word[startIndex:], nodeList)
			res = nodeList
		}
	}
	return res, ok
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
