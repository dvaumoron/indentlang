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
	"errors"
	"strconv"
	"strings"

	"github.com/dvaumoron/indentlang/types"
)

const AttributeName = "attribute"
const ListName = "List"

var ListId = types.NewIdentifier(ListName)

type stack[T any] struct {
	inner []T
}

func (s *stack[T]) push(e T) {
	s.inner = append(s.inner, e)
}

func (s *stack[T]) peek() T {
	return s.inner[len(s.inner)-1]
}

func (s *stack[T]) pop() T {
	last := len(s.inner) - 1
	res := s.inner[last]
	s.inner = s.inner[:last]
	return res
}

func Parse(str string) (*types.List, error) {
	var indentStack stack[int]
	indentStack.push(0)
	var listStack stack[*types.List]
	res := types.NewList()
	res.Add(ListId)
	listStack.push(res)
	lines := strings.Split(str, "\n")
LineLoop:
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" && trimmed[0] != '#' {
			var index int
			var char rune
			for index, char = range line {
				if char != ' ' && char != '\t' {
					if top := indentStack.peek(); top < index {
						indentStack.push(index)
						current := types.NewList()
						listStack.peek().Add(current)
						listStack.push(current)
					} else if top == index {
						listStack.pop()
						current := types.NewList()
						listStack.peek().Add(current)
						listStack.push(current)
					} else {
						indentStack.pop()
						listStack.pop()
						for top = indentStack.peek(); top > index; top = indentStack.peek() {
							indentStack.pop()
							listStack.pop()
						}
						if top < index {
							return nil, errors.New("identation not consistent")
						}
					}
					break
				}
			}
			words := make(chan string)
			done := make(chan types.NoneType)
			go handleWord(words, listStack, done)
			chars := make(chan rune)
			go sendChar(chars, line[index:])
			var buildingWord []rune
			for char := range chars {
				switch char {
				case ' ', '\t':
					buildingWord = sendReset(words, buildingWord)
				case '(', ')':
					buildingWord = sendReset(words, buildingWord)
					words <- string(char)
				case '"', '\'':
					var err error
					buildingWord, err = readUntil(buildingWord, char, chars, char)
					if err != nil {
						return nil, err
					}
				case '#':
					break LineLoop
				default:
					buildingWord = append(buildingWord, char)
				}
			}
			sendReset(words, buildingWord)
			close(words)
			<-done
		}
	}
	return res, nil
}

func handleWord(words <-chan string, listStack stack[*types.List], done chan<- types.NoneType) {
	for word := range words {
		if word == "(" {
			current := types.NewList()
			listStack.peek().Add(current)
			listStack.push(current)
		} else if word == ")" {
			listStack.pop()
		} else if handleCustomWord(word, listStack.peek()) {
			listStack.peek().Add(types.NewIdentifier(word))
		}
	}
	done <- types.None
}

var customRules = types.NewList()

func init() {
	customRules.Add(types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
		it := args.Iter()
		arg0, exist := it.Next()
		if exist {
			var str *types.String
			str, exist = arg0.(*types.String)
			if exist {
				if s := str.Inner; s == "true" {
					var arg1 types.Object
					arg1, exist = it.Next()
					if exist {
						var nodeList *types.List
						nodeList, exist = arg1.(*types.List)
						if exist {
							nodeList.Add(types.True)
						}
					}
				} else if s == "false" {
					var arg1 types.Object
					arg1, exist = it.Next()
					if exist {
						var nodeList *types.List
						nodeList, exist = arg1.(*types.List)
						if exist {
							nodeList.Add(types.False)
						}
					}
				} else if s == "None" {
					var arg1 types.Object
					arg1, exist = it.Next()
					if exist {
						var nodeList *types.List
						nodeList, exist = arg1.(*types.List)
						if exist {
							nodeList.Add(types.None)
						}
					}
				}
			}
		}
		return types.MakeBoolean(exist)
	}))
	customRules.Add(types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
		it := args.Iter()
		arg0, exist := it.Next()
		if exist {
			var str *types.String
			str, exist = arg0.(*types.String)
			if exist {
				s := str.Inner
				if s != ":=" && strings.Contains(s, ":") {
					var arg1 types.Object
					arg1, exist = it.Next()
					if exist {
						var nodeList *types.List
						nodeList, exist = arg1.(*types.List)
						if exist {
							nodeList2 := types.NewList()
							nodeList.Add(nodeList2)
							nodeList2.Add(ListId)
							for _, elem := range strings.Split(s, ":") {
								if elem == "" {
									nodeList2.Add(types.None)
								} else if handleCustomWord(elem, nodeList2) {
									nodeList2.Add(types.NewIdentifier(elem))
								}
							}
						}
					}
				}
			}
		}
		return types.MakeBoolean(exist)
	}))
	customRules.Add(types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
		it := args.Iter()
		arg0, exist := it.Next()
		if exist {
			var str *types.String
			str, exist = arg0.(*types.String)
			if exist {
				s := str.Inner
				i, err := strconv.ParseInt(s, 10, 64)
				if err == nil {
					var arg1 types.Object
					arg1, exist = it.Next()
					if exist {
						var nodeList *types.List
						nodeList, exist = arg1.(*types.List)
						if exist {
							nodeList.Add(types.NewInteger(i))
						}
					}
				} else {
					f, err := strconv.ParseFloat(s, 64)
					if err == nil {
						var arg1 types.Object
						arg1, exist = it.Next()
						if exist {
							var nodeList *types.List
							nodeList, exist = arg1.(*types.List)
							if exist {
								nodeList.Add(types.NewFloat(f))
							}
						}
					} else {
						exist = false
					}
				}
			}
		}
		return types.MakeBoolean(exist)
	}))
	customRules.Add(types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
		it := args.Iter()
		arg0, exist := it.Next()
		if exist {
			var str *types.String
			str, exist = arg0.(*types.String)
			if exist {
				if s := str.Inner; s[0] == '"' {
					var arg1 types.Object
					arg1, exist = it.Next()
					if exist {
						var nodeList *types.List
						nodeList, exist = arg1.(*types.List)
						if exist {
							s = strings.ReplaceAll(s[1:len(s)-1], "\\'", "'")
							nodeList.Add(types.NewString(s))
						}
					}
				}
			}
		}
		return types.MakeBoolean(exist)
	}))
	customRules.Add(types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
		it := args.Iter()
		arg0, exist := it.Next()
		if exist {
			var str *types.String
			str, exist = arg0.(*types.String)
			if exist {
				if s := str.Inner; s[0] == '\'' {
					var arg1 types.Object
					arg1, exist = it.Next()
					if exist {
						var nodeList *types.List
						nodeList, exist = arg1.(*types.List)
						if exist {
							skipApos := false
							size := len(s)
							s2 := make([]rune, 0, size)
							for _, char := range s[1 : size-1] {
								if skipApos {
									skipApos = false
									if char == '\'' {
										s2 = append(s2, char)
									} else {
										s2 = append(s2, '\\', char)
									}
								} else if char == '"' {
									s2 = append(s2, '\\', char)
								} else if char == '\\' {
									skipApos = true
								} else {
									s2 = append(s2, char)
								}
							}
							nodeList.Add(types.NewString(string(s2)))
						}
					}
				}
			}
		}
		return types.MakeBoolean(exist)
	}))
	customRules.Add(types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
		it := args.Iter()
		arg0, exist := it.Next()
		if exist {
			var str *types.String
			str, exist = arg0.(*types.String)
			if exist {
				if s := str.Inner; s[0] == '@' && len(s) > 1 {
					var arg1 types.Object
					arg1, exist = it.Next()
					if exist {
						var nodeList *types.List
						nodeList, exist = arg1.(*types.List)
						if exist {
							elems := strings.Split(s[1:], "=")
							attr := types.NewList()
							attr.AddCategory(AttributeName)
							attr.Add(types.NewString(elems[0]))
							if len(elems) > 1 {
								elem := elems[1]
								if handleCustomWord(elem, attr) {
									attr.Add(types.NewIdentifier(elem))
								}
							}
							nodeList.Add(attr)
						}
					}
				}
			}
		}
		return types.MakeBoolean(exist)
	}))
}

// Counter intuitive : return true when nothing have been done
func handleCustomWord(word string, list *types.List) bool {
	args := types.NewList()
	args.Add(types.NewString(word))
	args.Add(list)
	args.Add(customRules)
	res := !types.ForEach(customRules, func(object types.Object) bool {
		rule, success := object.(types.Appliable)
		if success {
			var boolean types.Boolean
			// The Apply must return true if it has created the node.
			boolean, success = rule.Apply(nil, args).(types.Boolean)
			success = success && boolean.Inner
		}
		return !success
	})
	return res
}

func sendChar(chars chan<- rune, line string) {
	for _, char := range line {
		chars <- char
	}
	close(chars)
}

func sendReset(words chan<- string, buildingWord []rune) []rune {
	if len(buildingWord) != 0 {
		words <- string(buildingWord)
		// doesn't realloc memmory
		buildingWord = buildingWord[:0]
	}
	return buildingWord
}

func readUntil(buildingWord []rune, char rune, chars <-chan rune, delim rune) ([]rune, error) {
	buildingWord = append(buildingWord, char)
	char, exist := <-chars
	for exist {
		buildingWord = append(buildingWord, char)
		if char == delim {
			break
		} else if char == '\\' {
			char, exist = <-chars
			if exist {
				buildingWord = append(buildingWord, char)
			} else {
				return nil, errors.New("unended string")
			}
		}
		char, exist = <-chars
	}
	return buildingWord, nil
}
