let pariEnCours;
let matchDisponible;
let registerModal;
let loginModal;
let loginBtn;
let registerBtn;
let registerForm;
let loginForm;
let idSession = "";
let pariEnCoursListe;
let matchDisponibleListe;

window.onload = function () {
    pariEnCours = document.getElementById("pariEnCours")
    matchDisponible = document.getElementById("matchDisponible")
    registerModal = document.getElementById("registerModal");
    loginModal = document.getElementById("loginModal");
    loginBtn = document.getElementById("loginBtn");
    registerBtn = document.getElementById("registerBtn");
    registerForm = document.getElementById("registerForm")
    loginForm = document.getElementById("loginForm")
    pariEnCoursListe = document.getElementById("pariEnCoursListe")
    matchDisponibleListe = document.getElementById("matchDisponibleListe")

    loginForm.addEventListener('submit', function (event) {
        event.preventDefault();
        let formData = new FormData(event.target);
        console.log(formData)
        fetch('/connexion', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept': 'application/json'
            },
            body: new URLSearchParams(formData)
        })
            .then(function (response) {
                return response.json();
            })
            .then(function (jsonData) {
                window.alert(jsonData["message"])
                if (jsonData["code"] === "200") {
                    idSession = jsonData["idSession"]
                    refreshMatchComing()
                }
            });
    })

    registerForm.addEventListener('submit', function (event) {
        event.preventDefault();
        let formData = new FormData(event.target);

        fetch('/user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept': 'application/json'
            },
            body: new URLSearchParams(formData)
        })
            .then(function (response) {
                return response.json();
            })
            .then(function (jsonData) {
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
    document.getElementById("loginRegister").value = ''
    document.getElementById("mailRegister").value = ''
    document.getElementById("passwordRegister").value = ''
}

function clearLoginForm() {
    document.getElementById("loginConnect").value = ''
    document.getElementById("passwordConnect").value = ''
}

function refreshMatchComing() {

    let params = new URLSearchParams()
    params.append("idSession", idSession)
    console.log(params.toString())
    fetch("/match?" + params.toString())
        .then(function (response) {
            return response.json();
        })
        .then(function (jsonData) {
            // window.alert(jsonData["message"])
            if (jsonData["code"] === "200") {
                let result = jsonData["result"];
                clearPariEnCoursListe()
                for (var i = 0; i < result.length; i++) {
                    var li = document.createElement('li')

                    var equipeAButton = document.createElement("button")
                    var equipeAText = document.createTextNode(result[i]["equipeA"])
                    equipeAButton.appendChild(equipeAText)

                    var coteContent = document.createTextNode(result[i]["cote"])
                    var equipeBButton = document.createElement("button")
                    var equipeBText = document.createTextNode(result[i]["equipeB"])
                    equipeBButton.appendChild(equipeBText)

                    var dateMatch = document.createTextNode((result[i]["date"]))

                    li.appendChild(equipeAButton)
                    li.appendChild(coteContent)
                    li.appendChild(equipeBButton)
                    li.appendChild(dateMatch)

                    matchDisponibleListe.append(li)
                }
            }
        });
}


window.onclick = function (event) {
    if (event.target === registerModal) {
        registerModal.style.display = "none";
        clearRegisterForm()
    } else if (event.target === loginModal) {
        loginModal.style.display = "none";
        clearLoginForm()
    }
}

function clearPariEnCoursListe() {
    document.getElementById("pariEnCoursListe").innerHTML = "";
}

function clearMatchAVenir() {
    document.getElementById("matchDisponibleListe").innerHTML = "";
}