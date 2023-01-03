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
	env      types.Environment
	waitings []chan<- types.Environment
	loaded   bool
}

// don't need mutex (only one coroutine access it)
var moduleCache = map[string]moduleCacheValue{}

func importer(requestReceiver <-chan importRequest, responseReceiver <-chan importResponse) {
	for {
		select {
		case request := <-requestReceiver:
			basePath := request.basePath
			totalPath := basePath + request.filePath
			value := moduleCache[totalPath]
			if value.loaded {
				responder := request.responder
				request.responder <- value.env
				close(responder)
			} else {
				waitings := value.waitings
				if len(waitings) == 0 {
					// nobody waiting, trying import
					go innerImporter(basePath, totalPath)
				}
				value.waitings = append(waitings, request.responder)
				moduleCache[totalPath] = value
			}
		case response := <-responseReceiver:
			path := response.path
			env := response.env
			// send the imported to all waitings
			for _, responder := range moduleCache[path].waitings {
				responder <- env
				close(responder)
			}
			// save the computed env & reset the list of waiting
			moduleCache[path] = moduleCacheValue{env: env, loaded: true}
		}
	}
}

const ImportName = "Import"

func innerImporter(basePath, totalPath string) {
	env := types.NewLocalEnvironment(Builtins)
	env.StoreStr(ImportName, makeCheckedImportDirective(basePath))
	// nested environment to isolate the directive Import, this avoid copying
	// (use types.Environment to get untyped nil)
	var local types.Environment = types.NewLocalEnvironment(env)
	tmplData, err := os.ReadFile(totalPath)
	if err == nil {
		var node types.Object
		node, err = parser.Parse(string(tmplData))
		if err == nil {
			node.Eval(local)
		} else {
			local = nil
		}
	} else {
		local = nil
	}
	responseToImporter <- importResponse{path: totalPath, env: local}
}

var directiveMutex sync.RWMutex
var directiveCache = map[string]types.NativeAppliable{}

func MakeImportDirective(basePath string) types.NativeAppliable {
	return makeCheckedImportDirective(CheckPath(basePath))
}

// internal version where basePath must end with a "/"
func makeCheckedImportDirective(basePath string) types.NativeAppliable {
	directiveMutex.RLock()
	res, ok := directiveCache[basePath]
	directiveMutex.RUnlock()
	if !ok {
		directiveMutex.Lock()
		res, ok = directiveCache[basePath]
		if !ok {
			res = types.MakeNativeAppliable(func(env types.Environment, itArgs types.Iterator) types.Object {
				arg0, _ := itArgs.Next()
				filePath, ok := arg0.Eval(env).(types.String)
				if ok {
					response := make(chan types.Environment)
					requestToImporter <- importRequest{
						basePath: basePath, filePath: string(filePath), responder: response,
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

// add an ending "/" if necessary
func CheckPath(path string) string {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path
}
