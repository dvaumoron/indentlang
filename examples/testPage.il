Import "main.il"

:= WidgetHeader (script @type="text/javascript" @src="/static/test.js")

:= WidgetBody
    List
        ul
            For ii (Range 3)
                := index (+ ii 1)
                li (a @href="/page?num=":ii "page num&eacute;ro " index)
        div (a @href="/otherPage" "Go there")
