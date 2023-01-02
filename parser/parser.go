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
const ListId types.Identifier = "List"

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

func newStack[T any]() *stack[T] {
	return &stack[T]{}
}

func Parse(str string) (*types.List, error) {
	indentStack := newStack[int]()
	indentStack.push(0)
	listStack := newStack[*types.List]()
	res := types.NewList(ListId)
	listStack.push(res)
	manageOpen(listStack)
	var err error
LineLoop:
	for _, line := range strings.Split(str, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" && trimmed[0] != '#' {
			index := 0
			var char rune
			for index, char = range line {
				if char != ' ' && char != '\t' {
					if top := indentStack.peek(); top < index {
						indentStack.push(index)
						manageOpen(listStack)
					} else if top == index {
						listStack.pop()
						manageOpen(listStack)
					} else {
						indentStack.pop()
						listStack.pop()
						for top = indentStack.peek(); top > index; top = indentStack.peek() {
							indentStack.pop()
							listStack.pop()
						}
						if top < index {
							err = errors.New("identation not consistent")
							break LineLoop
						}
						listStack.pop()
						manageOpen(listStack)
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
					buildingWord, err = readUntil(buildingWord, chars, char)
					if err != nil {
						break LineLoop
					}
				case '#':
					finishLine(words, buildingWord, done)
					continue LineLoop
				default:
					buildingWord = append(buildingWord, char)
				}
			}
			finishLine(words, buildingWord, done)
		}
	}
	return res, err
}

func manageOpen(listStack *stack[*types.List]) {
	current := types.NewList()
	listStack.peek().Add(current)
	listStack.push(current)
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
	unended := true
	buildingWord = append(buildingWord, delim)
CharLoop:
	for char := range chars {
		buildingWord = append(buildingWord, char)
		switch char {
		case delim:
			unended = false
			break CharLoop
		case '\\':
			buildingWord = append(buildingWord, <-chars)
		}
	}
	var err error
	if unended {
		err = errors.New("unended string")
	}
	return buildingWord, err
}

func finishLine(words chan<- string, buildingWord []rune, done <-chan types.NoneType) {
	sendReset(words, buildingWord)
	close(words)
	<-done
}
