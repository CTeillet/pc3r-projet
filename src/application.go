package main

import (
	"gitlab.com/CTeillet/pc3r-projet/src/bet"
	"gitlab.com/CTeillet/pc3r-projet/src/coins"
	"gitlab.com/CTeillet/pc3r-projet/src/connexion"
	"gitlab.com/CTeillet/pc3r-projet/src/match"
	"gitlab.com/CTeillet/pc3r-projet/src/message"
	"gitlab.com/CTeillet/pc3r-projet/src/user"
	"gitlab.com/CTeillet/pc3r-projet/src/utils"
	"log"
	"net/http"
	"os"
	"time"
)

func handleUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Printf("User\n")
	if err != nil {
		handleProblem(w, r)
	}
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
	err := r.ParseForm()
	log.Printf("Bet\n")
	if err != nil {
		handleProblem(w, r)
	}
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
	err := r.ParseForm()
	log.Printf("Match\n")
	if err != nil {
		handleProblem(w, r)
	}
	switch r.Method {
	case "GET":
		match.GetMatch(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleConnexion(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Printf("Connexion\n")
	if err != nil {
		handleProblem(w, r)
	}
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
	err := r.ParseForm()
	log.Printf("Coins\n")
	if err != nil {
		handleProblem(w, r)
	}
	switch r.Method {
	case "POST":
		coins.Generate(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Printf("Message\n")
	if err != nil {
		handleProblem(w, r)
	}
	switch r.Method {
	case "GET":
		message.GetMessage(w, r)
	case "POST":
		message.PostMessage(w, r)
	case "DELETE":
		message.DeleteMessage(w, r)
	default:
		handleProblem(w, r)
	}
}

//HandleHome
func handleHome(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Welcome\n")
	utils.SendResponse(w, http.StatusOK, `{"message":"hello world!"}`)
}

func handleProblem(w http.ResponseWriter, _ *http.Request) {
	log.Printf("A problem appear\n")
	utils.SendResponse(w, http.StatusInternalServerError, `{"message":"problem"}`)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	f, _ := os.Create("/var/log/golang/golang-server.log")

	log.SetOutput(f)

	updateComingMatches()
	updateResultMatchesAndBet()

	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("/web"))))

	http.HandleFunc("/", handleHome)

	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/bet", handleBet)
	http.HandleFunc("/match", handleMatch)
	http.HandleFunc("/connexion", handleConnexion)
	http.HandleFunc("/coins", handleCoins)
	http.HandleFunc("/message", handleMessage)

	log.Printf("Listening on port %s\n\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err.Error())
	}

	err = f.Close()
	if err != nil {
		panic(err.Error())
	}
}

func updateComingMatches() {
	ticker := time.NewTicker(24 * time.Hour)
	match.LoadComingMatchFor2Week()
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				match.LoadComingMatchFor2Week()
			}
		}
	}()
}

func updateResultMatchesAndBet() {
	ticker := time.NewTicker(1 * time.Hour)
	match.LoadResultMatchFor3Hours()
	bet.UpdateResult1Hour()
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				match.LoadResultMatchFor3Hours()
				bet.UpdateResult1Hour()
			}
		}
	}()
}
