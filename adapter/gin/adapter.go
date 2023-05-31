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

package ginadapter

import (
	"net/http"

	"github.com/dvaumoron/indentlang/adapter"
	"github.com/dvaumoron/indentlang/template"
	"github.com/gin-gonic/gin/render"
)

// match Render interface from gin.
type indentlangHTML struct {
	Template template.Template
	Data     any
}

func (r indentlangHTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return r.Template.Execute(w, r.Data)
}

const contentTypeName = "Content-Type"

var htmlContentType = []string{"text/html; charset=utf-8"}

// Writes HTML ContentType.
func (r indentlangHTML) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header[contentTypeName]; len(val) == 0 {
		header[contentTypeName] = htmlContentType
	}
}

// match HTMLRender interface from gin.
type indentlangHTMLRender struct {
	Templates map[string]template.Template
}

func (r indentlangHTMLRender) Instance(name string, data any) render.Render {
	return indentlangHTML{
		Template: r.Templates[name],
		Data:     data,
	}
}

// Use this method to init the HTMLRender in a gin Engine.
func LoadTemplatesAsRender(templatesPath string) (render.HTMLRender, error) {
	templates, err := adapter.LoadTemplates(templatesPath)
	if err != nil {
		return nil, err
	}
	return indentlangHTMLRender{Templates: templates}, nil
}
