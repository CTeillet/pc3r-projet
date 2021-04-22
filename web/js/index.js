let registerModal;
let loginModal;
let loginBtn;
let registerBtn;
let registerForm;
let loginForm;
let idSession = "";
let pariEnCoursListe;
let matchDisponibleListe;
let matchDisponibleButton;
// let parisFinisListe;

window.onload = function () {
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
                for (let i = 0; i < result.length; i++) {
                    const form = document.createElement('form')
                    form.name=result[i]["id"]

                    const submitButton = document.createElement('button')

                    submitButton.type="submit"

                    submitButton.appendChild(document.createTextNode("Soumettre Pari"))

                    const li = document.createElement('li');

                    const ul = document.createElement('ul');

                    const sport = document.createTextNode("Sport : " + result[i]["sport"]);
                    const sportLi = document.createElement('li');
                    sportLi.appendChild(sport)

                    const league = document.createTextNode("League : " + result[i]["league"]);
                    const leagueLi = document.createElement('li');
                    leagueLi.appendChild(league)

                    const date = document.createTextNode('Date : ' + result[i]["date"]);
                    const dateLi = document.createElement('li');
                    dateLi.appendChild(date)

                    const equipe = document.createTextNode(result[i]["equipeA"] + " vs " + result[i]["equipeB"]);
                    const equipeLi = document.createElement('li');
                    equipeLi.appendChild(equipe)

                    const cote = document.createTextNode("Cote : " + result[i]["cote"]);
                    const coteLi = document.createElement('li');
                    coteLi.appendChild(cote)

                    const montant = document.createElement('input');
                    const montantLi = document.createElement('li');
                    const montantTxt = document.createTextNode("Montant ");
                    montantLi.append(montantTxt, montant)
                    montant.type= 'number'
                    montant.value=0
                    montant.min=0

                    const vainqueurLi = document.createElement('li');

                    const equipeARadio = document.createElement('input');
                    equipeARadio.type='radio'
                    equipeARadio.value=result[i]["equipeA"]
                    equipeARadio.name='vainqueur'+result[i]["id"]
                    equipeARadio.id='equipeA'+result[i]["id"]

                    const equipeALabel = document.createElement('label');
                    equipeALabel.htmlFor='equipeA'+result[i]["id"]

                    const equipeALabelText = document.createTextNode(result[i]["equipeA"]);

                    equipeALabel.appendChild(equipeALabelText)

                    const equipeBRadio = document.createElement('input');
                    equipeBRadio.type='radio'
                    equipeBRadio.value=result[i]["equipeB"]
                    equipeBRadio.name='vainqueur'+result[i]["id"]
                    equipeBRadio.id='equipeB'+result[i]["id"]

                    const equipeBLabel = document.createElement('label');
                    equipeALabel.htmlFor='equipeB'+result[i]["id"]

                    const equipeBLabelText = document.createTextNode(result[i]["equipeB"]);

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

                    form.appendChild(ul)
                    form.appendChild(submitButton)

                    li.appendChild(form)

                    matchDisponibleListe.append(li)

                    submitButton.onclick=function (event){
                        event.preventDefault()
                        let idMatch = event.target.form.name
                        let vainqueur = document.querySelector('input[name="vainqueur'+idMatch+'"]:checked').value;
                        if (montant.value!==0 && vainqueur!=="" ){
                            params = new URLSearchParams()
                            console.log("Montant : " + montant.value)
                            console.log("Cote : " + result[i]["cote"])
                            params.append("idSession", idSession)
                            params.append("idMatch", idMatch)
                            params.append("equipeGagnante", vainqueur)
                            params.append("cote", result[i]["cote"])
                            params.append("montant", montant.value)

                            fetch("/bet", {
                                method: 'POST',
                                headers: {
                                    'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                                    'Accept': 'application/json'
                                },
                                body: params
                            })
                                .then(function (response) {
                                    console.log(response)
                                    return response.json()
                                })
                                .then(function (jsonData) {
                                    window.alert(jsonData["message"])
                                })
                        }
                    }

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

function refreshBetHistory() {

}