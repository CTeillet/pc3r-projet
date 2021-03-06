package utils

import (
	"gitlab.com/CTeillet/pc3r-projet/src/database"
	"net/http"
	"strconv"
	"time"
)

func SendResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	s := `{"code":"` + strconv.Itoa(status) + `",` + message[1:]
	_, err := w.Write([]byte(s))
	if err != nil {
		panic(err.Error())
	}
}

func IsUser(login string, password string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	count := 0
	err := db.QueryRow("Select count(*) From Utilisateur where login=? and password=?;", login, password).Scan(&count)
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}
	if count == 1 {
		return true
	}
	return false
}

type Connexion struct {
	id        int
	login     string
	idSession string
	date      time.Time
}

func IsConnectedIdSession(idSession string) string {
	db := database.Connect()
	if db == nil {
		return ""
	}
	c := Connexion{}
	err := db.QueryRow("Select * From Connexion where idSession=?;", idSession).Scan(&c.id, &c.login, &c.idSession, &c.date)
	if err != nil {
		return ""
	}
	err = db.Close()
	if err != nil {
		return ""
	}

	t := time.Now().UTC()
	comp := c.date.Add(15 * time.Minute)

	if comp.After(t) {
		majConnexion(idSession)
		return c.login
	} else {
		RemoveConnection(idSession)
		return ""
	}
}

func IsConnectedLogin(login string) string {
	db := database.Connect()
	if db == nil {
		return ""
	}
	idSession := ""
	err := db.QueryRow("Select idSession From Connexion where login=?", login).Scan(&idSession)

	if err != nil || idSession == "" {
		return ""
	}
	return idSession

}

//getLogin : plus utilisée car le login était dans c
func _(idSession string) string {
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

func RemoveConnection(idSession string) bool {
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
