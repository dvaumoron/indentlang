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

var customRules = types.NewList()

func init() {
	customRules.Add(types.MakeNativeAppliable(constantRule))
	customRules.Add(types.MakeNativeAppliable(listLiteralRule))
	customRules.Add(types.MakeNativeAppliable(numberRule))
	customRules.Add(types.MakeNativeAppliable(stringLiteralRule))
	customRules.Add(types.MakeNativeAppliable(stringLiteralRule2))
	customRules.Add(types.MakeNativeAppliable(attributeRule))
}

func constantRule(env types.Environment, args types.Iterable) types.Object {
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
}

func listLiteralRule(env types.Environment, args types.Iterable) types.Object {
	it := args.Iter()
	arg0, exist := it.Next()
	if exist {
		var str *types.String
		str, exist = arg0.(*types.String)
		if exist {
			s := str.Inner
			if s != SetName && strings.Contains(s, ":") {
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
}

func numberRule(env types.Environment, args types.Iterable) types.Object {
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
}

func stringLiteralRule(env types.Environment, args types.Iterable) types.Object {
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
}

func stringLiteralRule2(env types.Environment, args types.Iterable) types.Object {
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
}

func attributeRule(env types.Environment, args types.Iterable) types.Object {
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
}
