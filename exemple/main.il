html
    head
        meta @charset="utf-8"
        title "Title - " PageTitle
        link @rel="stylesheet" @type="text/css" @href="/static/main.css"
        script @type="text/javascript" @src="/static/main.js"
        Js
    body
        header
            h1 PageTitle
            if LoginUrl
                p
                    if Login (Welcome "&nbsp;" Login "&nbsp;")
                    a @href=LoginUrl LoginLinkName
            if ErrorMsg
                div @class="error" ErrorMsg
            p
                for page Ariane
                    := url (. page Url)
                    "/&nbsp;" (a @href=url (. page Name)) "&nbsp;"
            ul
                for page SubPages
                    := url (. page Url)
                    li (a @href=url (. page Name))
        main
            Body
        footer
