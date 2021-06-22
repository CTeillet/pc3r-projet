let idSession = sessionStorage.getItem('idSession');
// console.log("ID SESSION " + idSession)

if(idSession!=="" && idSession!==null){
    let a = new Twitch.Embed("twitch-embed", {
        width: "100%",
        height: 1000,
        channel: "otplol_",
        autoplay: true
    });

    setInterval(function (param) {
        let params = new URLSearchParams()
        params.append("idSession", idSession)
        params.append('montant', 50)
        fetch('/coins', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept': 'application/json'
            },
            body: new URLSearchParams(params)
        })
            .then(function (response) {
                return response.json();
            })
            .then(function (jsonData) {
                // console.log(jsonData);
                window.alert(jsonData["message"])
            });
    }, 50000)
}
