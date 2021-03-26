package main

import (
	"gitlab.com/CTeillet/pc3r-projet/bet"
	"gitlab.com/CTeillet/pc3r-projet/coins"
	"gitlab.com/CTeillet/pc3r-projet/connexion"
	"gitlab.com/CTeillet/pc3r-projet/match"
	"gitlab.com/CTeillet/pc3r-projet/message"
	"gitlab.com/CTeillet/pc3r-projet/user"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	id       int
	login    string
	password string
	mail     string
}

type Bet struct {
	id             int
	idMatch        int
	equipeGagnante string
	cote           float32
	montant        int
}

type Match struct {
	id      int
	equipeA string
	equipeB string
	cote    float32
}

type Connexion struct {
	id        int
	login     string
	idSession string
	date      time.Time
}

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
		handleProblem(w, r)
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
		handleProblem(w, r)
	}
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		match.GetMatch(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleConnexion(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		connexion.Connect(w, r)
	case "DELETE":
		connexion.Disconnect(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleCoins(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		coins.Generate(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		message.GetMessage(w, r)
	case "POST":
		message.PostMessage(w, r)
	case "DELETE":
		message.DeleteMessage(w, r)
	}
}

func handleHome(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message":"hello world!"}`))
}

func handleProblem(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message":"problem"}`))
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
	http.HandleFunc("/messsage", handleMessage)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
