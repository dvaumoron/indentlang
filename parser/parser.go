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
					buildingWord = append(buildingWord, char)
					var err error
					buildingWord, err = readUntil(buildingWord, chars, char)
					if err != nil {
						return nil, err
					}
				case '#':
					finishLine(words, buildingWord, done)
					break LineLoop
				default:
					buildingWord = append(buildingWord, char)
				}
			}
			finishLine(words, buildingWord, done)
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

// an empty environment to execute custom rules
var BuiltinsCopy types.Environment = types.MakeBaseEnvironment()

// Counter intuitive : return true when nothing have been done
func handleCustomWord(word string, list *types.List) bool {
	args := types.NewList()
	args.Add(types.NewString(word))
	args.Add(list)
	args.Add(customRules)
	res := true
	types.ForEach(customRules, func(object types.Object) bool {
		rule, ok := object.(types.Appliable)
		if ok {
			// The Apply must return true if it has created the node.
			boolean, ok := rule.Apply(BuiltinsCopy, args).(types.Boolean)
			res = !(ok && boolean.Inner)
		}
		return res
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

func readUntil(buildingWord []rune, chars <-chan rune, delim rune) ([]rune, error) {
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

func finishLine(words chan<- string, buildingWord []rune, done <-chan types.NoneType) {
	sendReset(words, buildingWord)
	close(words)
	<-done
}
