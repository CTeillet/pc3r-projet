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
var matchDisponibleButton

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
    matchDisponibleButton = document.getElementById("matchDisponibleButton")

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
                    refreshBet()
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
            if (jsonData["code"] === "200") {
                let result = jsonData["result"];
                clearMatchAVenir()
                for (var i = 0; i < result.length; i++) {
                    var li = document.createElement('li')

                    var ul = document.createElement('ul')

                    var sport = document.createTextNode("Sport : " + result[i]["sport"])
                    var sportLi = document.createElement('li')
                    sportLi.appendChild(sport)

                    var league = document.createTextNode("League : " + result[i]["league"])
                    var leagueLi = document.createElement('li')
                    leagueLi.appendChild(league)

                    var date = document.createTextNode('Date : ' + result[i]["date"])
                    var dateLi = document.createElement('li')
                    dateLi.appendChild(date)

                    var equipe = document.createTextNode(result[i]["equipeA"] + " vs " + result[i]["equipeB"])
                    var equipeLi = document.createElement('li')
                    equipeLi.appendChild(equipe)

                    var cote = document.createTextNode("Cote : " +  result[i]["cote"])
                    var coteLi = document.createElement('li')
                    coteLi.appendChild(cote)

                    var montant = document.createElement('input')
                    var montantLi = document.createElement('li')
                    var montantTxt = document.createTextNode("Montant ")
                    montantLi.append(montantTxt, montant)
                    montant.type= 'number'

                    var vainqueurLi = document.createElement('li')

                    var equipeARadio = document.createElement('input')
                    equipeARadio.type='radio'
                    equipeARadio.value=result[i]["equipeA"]
                    equipeARadio.name='vainqueur'+result[i]["id"]
                    equipeARadio.id='equipeA'+result[i]["id"]

                    var equipeALabel = document.createElement('label')
                    equipeALabel.htmlFor='equipeA'+result[i]["id"]

                    var equipeALabelText = document.createTextNode(result[i]["equipeA"])

                    equipeALabel.appendChild(equipeALabelText)

                    var equipeBRadio = document.createElement('input')
                    equipeBRadio.type='radio'
                    equipeBRadio.value=result[i]["equipeB"]
                    equipeBRadio.name='vainqueur'+result[i]["id"]
                    equipeBRadio.id='equipeB'+result[i]["id"]

                    var equipeBLabel = document.createElement('label')
                    equipeALabel.htmlFor='equipeB'+result[i]["id"]

                    var equipeBLabelText = document.createTextNode(result[i]["equipeB"])

                    equipeBLabel.appendChild(equipeBLabelText)

                    vainqueurLi.appendChild(equipeARadio)
                    vainqueurLi.appendChild(equipeALabel)
                    vainqueurLi.appendChild(equipeBRadio)
                    vainqueurLi.appendChild(equipeBLabel)

                    ul.appendChild(sportLi)
                    ul.appendChild(leagueLi)
                    ul.appendChild(dateLi)
                    ul.appendChild(equipeLi)
                    ul.appendChild(coteLi)
                    ul.appendChild(vainqueurLi)
                    ul.appendChild(montantLi)

                    li.appendChild(ul)
                    matchDisponibleListe.append(li)

                }
            }
        });
}

function refreshBet() {
    clearPariEnCoursListe()
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