var pariEnCours
var matchDisponible
var registerModal;
var loginModal;
var loginBtn;
var registerBtn;

window.onload = function () {
    pariEnCours = document.getElementById("pariEnCours")
    matchDisponible = document.getElementById("matchDisponible")
    registerModal = document.getElementById("registerModal");
    loginModal = document.getElementById("loginModal");
    loginBtn = document.getElementById("loginBtn");
    registerBtn = document.getElementById("registerBtn");

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