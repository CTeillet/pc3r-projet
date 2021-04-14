var pariEnCours
var matchDisponible
var registerModal;
var loginModal;
var loginBtn;
var registerBtn;
var registerForm;
var idSession ="";

window.onload = function () {
    pariEnCours = document.getElementById("pariEnCours")
    matchDisponible = document.getElementById("matchDisponible")
    registerModal = document.getElementById("registerModal");
    loginModal = document.getElementById("loginModal");
    loginBtn = document.getElementById("loginBtn");
    registerBtn = document.getElementById("registerBtn");
    registerForm = document.getElementById("registerForm")

    registerForm.addEventListener('submit', function(event) {
        event.preventDefault();

        const formData = new FormData(event.target);

        var responsePromise = fetch('/user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept':'application/json'
            },
            body: new URLSearchParams(formData)
        })
        responsePromise
            .then(function(response) {
                return response.json();
            })
            .then(function(jsonData) {
                console.log(jsonData);
                window.alert(jsonData["message"])
            });
    })

    loginBtn.onclick = function () {
        loginModal.style.display = "block";
    }

    registerBtn.onclick = function () {
        registerModal.style.display = "block";
    }
}

function clearRegisterForm() {
    document.getElementById("loginRegister").value=''
    document.getElementById("mailRegister").value=''
    document.getElementById("passwordRegister").value=''
}

function clearLoginForm() {
    document.getElementById("loginConnect").value=''
    document.getElementById("passwordConnect").value=''
}

window.onclick = function(event) {
    if (event.target === registerModal) {
        registerModal.style.display = "none";
        clearRegisterForm()
    }
    else if (event.target === loginModal){
        loginModal.style.display = "none";
        clearLoginForm()
    }
}