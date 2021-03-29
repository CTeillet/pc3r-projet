package connexion

import (
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/user"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"math/rand"
	"net/http"
	"time"
)

type Connexion struct {
	id        int
	login     string
	idSession string
	date      time.Time
}

func Connect(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if user.IsUser(login, password) {
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

	if IsConnected(idSession) != "" && idSession != "" {
		if removeConnection(idSession) {
			utils.SendResponse(w, http.StatusOK, `{"message":"user disconnected"}`)
		} else {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message":"a problem appeared"}`)
		}
	} else {
		utils.SendResponse(w, http.StatusForbidden, `{"message":"user was not connected"}`)
	}
}

func IsConnected(idSession string) string {
	db := database.Connect()
	if db == nil {
		return ""
	}
	c := Connexion{}
	err := db.QueryRow("Select * From Connexion where idSession=?;", idSession).Scan(&c.id, &c.login, &c.idSession, c.date)
	if err != nil {
		return ""
	}
	err = db.Close()
	if err != nil {
		return ""
	}

	t := time.Now()
	comp := c.date.Add(15 * time.Minute)

	if comp.After(t) {
		majConnexion(idSession)
		login := getLogin(idSession)
		return login
	} else {
		removeConnection(idSession)
		return ""
	}

}

func getLogin(idSession string) string {
	db := database.Connect()
	if db == nil {
		return ""
	}
	err := db.Close()
	if err != nil {
		return ""
	}
	login := ""

	err = db.QueryRow("Select login From Connexion where idSession=?", idSession).Scan(&login)
	if err != nil {
		return ""
	}

	return login
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

func removeConnection(idSession string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	res, err := db.Exec("Delete from Connexion where idSession=?;", idSession)
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}

	r, _ := res.RowsAffected()
	if r == 0 {
		return false
	} else {
		return true
	}
}

func majConnexion(idSession string) {
	db := database.Connect()
	_, err := db.Exec("UPDATE Connexion set date=now() where idSession=?;", idSession)
	if err != nil {
		return
	}
	err = db.Close()
	if err != nil {
		return
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
