package connexion

import (
	"gitlab.com/CTeillet/pc3r-projet/src/database"
	"gitlab.com/CTeillet/pc3r-projet/src/utils"
	"math/rand"
	"net/http"
)

func Connect(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if utils.IsUser(login, password) {
		idSession := addConnexion(login)
		if idSession != "" {
			utils.SendResponse(w, http.StatusOK, `{"message":"user connected", "idSession":"`+idSession+`"}`)
		} else {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message":"problem with database"}`)
		}

	} else {
		utils.SendResponse(w, http.StatusForbidden, `{"message":"unknown user"}`)
	}
}

func Disconnect(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")

	if utils.IsConnected(idSession) != "" && idSession != "" {
		if utils.RemoveConnection(idSession) {
			utils.SendResponse(w, http.StatusOK, `{"message":"user disconnected"}`)
		} else {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message":"a problem appeared"}`)
		}
	} else {
		utils.SendResponse(w, http.StatusForbidden, `{"message":"user was not connected"}`)
	}
}

func addConnexion(login string) string {
	db := database.Connect()
	if db == nil {
		return ""
	}
	idSession := randSeq()
	res, err := db.Exec("INSERT INTO Connexion (login, idSession) VALUES (?, ?);", login, idSession)
	if err != nil {
		return ""
	}
	err = db.Close()
	if err != nil {
		return ""
	}

	r, _ := res.RowsAffected()
	if r == 0 {
		return ""
	} else {
		return idSession
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func randSeq() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
