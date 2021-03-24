package main

import (
	"gitlab.com/CTeillet/pc3r-projet/bet"
	"gitlab.com/CTeillet/pc3r-projet/coins"
	"gitlab.com/CTeillet/pc3r-projet/connexion"
	"gitlab.com/CTeillet/pc3r-projet/match"
	"gitlab.com/CTeillet/pc3r-projet/user"
	"net/http"
	"os"
)

//type Personne struct {
//
//}
//
//type Pari struct {
//
//}
//
//type Match struct {
//
//}

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user.GetUser(w, r)
	case "POST":
		user.AddUser(w, r)
	case "PUT":
		user.ModifyUser(w, r)
	case "DELETE":
		user.DeleteUser(w, r)
	default:

	}
}

func handleBet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		bet.GetBet(w, r)
	case "POST":
		bet.AddBet(w, r)
	case "DELETE":
		bet.DeleteBet(w, r)
	default:

	}
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		match.GetMatch(w, r)
	default:

	}
}

func handleConnexion(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		connexion.Connect(w, r)
	case "DELETE":
		connexion.Disconnect(w, r)
	default:

	}
}

func handleCoins(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		coins.Generate(w, r)
	default:

	}
}

func handleHome(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message":"hello world!"}`))

}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/bet", handleBet)
	http.HandleFunc("/match", handleMatch)
	http.HandleFunc("/connexion", handleConnexion)
	http.HandleFunc("/coins", handleCoins)
	_ = http.ListenAndServe(":"+port, nil)

}
