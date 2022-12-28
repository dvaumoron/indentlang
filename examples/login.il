Import "main.il"

:= WidgetHeader (script @type="text/javascript" @src="/static/login.js")

:= WidgetBody
    form @action=CurrentUrl:"submit" @method="POST"
        div (label @for="Login" LoginLabel)
        div (input @type="text" @name="Login")
        div (label @for="Password" PasswordLabel)
        div (input @type="password" @name="Password")
        div @class="hidden"
            input @type="text" @name="Register"
            input @type="text" @name="Redirect" @value=Redirect
            input @type="text" @name="PrevUrlWithError" @value=PrevUrlWithError
        div (input @type="submit" @value=LoginLinkName)
        div (a @href="todo" RegisterLinkName)
