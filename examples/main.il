html
    head
        meta @charset="utf-8"
        title "Title - " Title
        link @rel="stylesheet" @type="text/css" @href="/static/main.css"
        script @type="text/javascript" @src="/static/main.js" ""
        WidgetHeader
    body
        header
            h1 Title
            p
                If UserName ("Welcome&nbsp;" UserName "&nbsp;")
                a @href=Link LinkLabel
        main
            Eval WidgetBody
        footer
            For arg Footers
                p arg
