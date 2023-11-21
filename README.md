# IndentLang

A templating language mainly to output balised language (html, xml, etc.).

## usage

Start by importing the template package :

```Go
import "github.com/dvaumoron/indentlang/template"
```

And use it in two step:

```Go
// parse a template
tmpl, err := template.ParsePath(tmplPath)
// and use it
err = tmpl.Execute(writer, data)
```

With the input (indentation matters):

```
html
    head
        meta @charset="utf-8"
        title "Hello World"
    body
        h1 @class="greetings" "Hello World"

```

The output will look like (cleaned):

```html
<html>
    <head>
        <meta charset="utf-8"/>
        <title>Hello World</title>
    </head>
    <body>
        <h1 class="greetings">Hello World</h1>
    </body>
</html>
```

The file [indentlang.go](indentlang.go) is an adapted copy of [engine.go](https://github.com/dvaumoron/ste/blob/master/engine.go) for demo and testing purpose (see [examples](examples)).

More examples can be found [here](https://github.com/dvaumoron/puzzletest/tree/main/templatedata/templates/indentlang).

See [API Documentation](https://pkg.go.dev/github.com/dvaumoron/indentlang/template).
