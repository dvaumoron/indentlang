Import "main"

:= WidgetHeader (script @type="text/javascript" @src="/static/test.js" "") # the empty string avoid self-closing

:= WidgetBody
    Quote
        List
            ul
                For i (Range (Integer Number))
                    := index (+ i 1)
                    li (a @href="/page?num=":i "page num&eacute;ro " index)
            div (a @href="/otherPage" "Go there")
