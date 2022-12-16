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
package builtins

import (
	"os"
	"sync"

	"github.com/dvaumoron/indentlang/parser"
	"github.com/dvaumoron/indentlang/types"
)

type importRequest struct {
	basePath  string
	filePath  string
	responder chan<- types.Environment
}

var requestToImporter chan<- importRequest

type importResponse struct {
	path string
	env  types.Environment
}

var responseToImporter chan<- importResponse

func init() {
	importChannel := make(chan importRequest)
	importChannel2 := make(chan importResponse)
	go importer(importChannel, importChannel2)
	requestToImporter = importChannel
	responseToImporter = importChannel2
}

type moduleCacheValue struct {
	env     types.Environment
	waiting []chan<- types.Environment
}

var moduleCache = map[string]moduleCacheValue{}

func importer(requestReceiver <-chan importRequest, responseReceiver <-chan importResponse) {
	for {
		select {
		case request := <-requestReceiver:
			basePath := request.basePath
			totalPath := basePath + request.filePath
			value := moduleCache[totalPath]
			if env := value.env; env == nil {
				list := value.waiting
				if len(list) == 0 {
					// import non existing
					go innerImporter(basePath, totalPath)
				}
				value.waiting = append(list, request.responder)
				moduleCache[totalPath] = value
			} else {
				request.responder <- env
			}
		case response := <-responseReceiver:
			path := response.path
			env := response.env
			for _, responder := range moduleCache[path].waiting {
				responder <- env
			}
			// save the computed env & reset the list of waiting
			moduleCache[path] = moduleCacheValue{env: env}
		}
	}
}

func innerImporter(basePath, totalPath string) {
	env := types.NewLocalEnvironment(Builtins)
	env.StoreStr("Import", MakeImportDirective(basePath))
	tmplData, err := os.ReadFile(totalPath)
	if err == nil {
		var node types.Object
		node, err = parser.Parse(string(tmplData))
		if err == nil {
			node.Eval(env)
		} else {
			env = nil
		}
	} else {
		env = nil
	}
	responseToImporter <- importResponse{path: totalPath, env: env}
}

var directiveMutex sync.RWMutex
var directiveCache = map[string]types.NativeAppliable{}

func MakeImportDirective(basePath string) types.NativeAppliable {
	directiveMutex.RLock()
	res, exist := directiveCache[basePath]
	directiveMutex.RUnlock()
	if !exist {
		directiveMutex.Lock()
		res, exist = directiveCache[basePath]
		if !exist {
			res = types.MakeNativeAppliable(func(env types.Environment, args *types.List) types.Object {
				filePath, success := args.Load(types.NewInteger(0)).(*types.String)
				if success {
					response := make(chan types.Environment)
					requestToImporter <- importRequest{
						basePath: basePath, filePath: filePath.Inner, responder: response,
					}
					otherEnv := <-response
					if otherEnv != nil {
						otherEnv.CopyTo(env)
					}
				}
				return types.None
			})
			directiveCache[basePath] = res
		}
		directiveMutex.Unlock()
	}
	return res
}
