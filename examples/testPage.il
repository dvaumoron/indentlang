Import "main.il"

:= WidgetHeader (script @type="text/javascript" @src="/static/test.js" "")

Func WidgetBody (n)
    Return
        ul
            For ii (Range n)
                := index (+ ii 1)
                li (a @href="/page?num=":ii "page num&eacute;ro " index)
        div (a @href="/otherPage" "Go there")
