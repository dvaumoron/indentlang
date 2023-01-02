html
    head
        meta @charset="utf-8"
        title "Title - " value1
        link @rel="stylesheet" @type="text/css" @href="/static/main.css"
        script @type="text/javascript" @src="/static/main.js" ""
        WidgetHeader
    body
        header
            h1 value1
            p
                If value2 ("Welcome&nbsp;" value2 "&nbsp;")
                a @href=value3 value4
        main
            WidgetBody 5
        footer
            For arg values
                p arg
