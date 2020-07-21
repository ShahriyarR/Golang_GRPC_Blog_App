import Navigo from 'navigo'

const router = new Navigo()

router
    .on("/", function () {
        document.body.innerHTML = "Home"
    })
    .on("/login", function () {
        document.body.innerHTML = ""
        const loginDiv = document.createElement('div')
        loginDiv.classList.add("login-div")

        const loginLabel = document.createElement('h1')
        loginLabel.innerText = "Login"
        loginDiv.appendChild(loginLabel)

        const loginForm = document.createElement('form')

        const loginInputLabel = document.createElement("label")
        loginInputLabel.innerText = "Username / E-Mail"
        loginInputLabel.setAttribute("for", "login-input")
        loginDiv.append(loginInputLabel)
        // TODO: start to watch Episode 5 from 19:43
        const loginInput = document.createElement('input')
        loginInput.id = "login-input"
        loginInput.setAttribute("type", "text")
        loginInput.setAttribute("placeholder", "John42")
        loginForm.appendChild(loginInput)
        loginDiv.appendChild(loginForm)

        document.body.appendChild(loginDiv)
    })
    .resolve()