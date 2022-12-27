html
    head
        meta @charset="utf-8"
        title "Title - " PageTitle
        link @rel="stylesheet" @type="text/css" @href="/static/main.css"
        script @type="text/javascript" @src="/static/main.js"
        WidgetHeader
    body
        header
            h1 PageTitle
            If LoginUrl
                p
                    If Login (Welcome "&nbsp;" Login "&nbsp;")
                    a @href=LoginUrl LoginLinkName
            If ErrorMsg
                div @class="error" ErrorMsg
            p
                For page Ariane
                    := url (. page Url)
                    "/&nbsp;" (a @href=url (. page Name)) "&nbsp;"
            ul
                For page SubPages
                    := url (. page Url)
                    li (a @href=url (. page Name))
        main
            WidgetBody
        footer
