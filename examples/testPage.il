Import "main.il"

:= WidgetHeader (script @type="text/javascript" @src="/static/test.js" "") # the empty string avoid self-closing

:= WidgetBody
    Quote
        List
            ul
                For ii (Range (Integer value5))
                    := index (+ ii 1)
                    li (a @href="/page?num=":ii "page num&eacute;ro " index)
            div (a @href="/otherPage" "Go there")
