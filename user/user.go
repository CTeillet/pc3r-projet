package user

import (
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"net/http"
)

type User struct {
	login    string
	password string
	mail     string
	cagnotte int
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	idSession := req.FormValue("idSession")
	login := res.FormValue("login")
	if connexion.IsConnected(idSession) {
		user := searchUser(login)
		if user != nil {
			if user.login != nil {
				utils.SendResponse(res, http.StatusOK, `{"message":"user found", "login":"`+user.login+`", "mail":"`+user.mail+`", "cagnotte":"`+user.cagnotte+`"}`)
			} else {
				utils.SendResponse(res, http.StatusForbidden, `{"message":"problem login user don't exist"}`)
			}
		} else {
			utils.SendResponse(res, http.StatusInternalServerError, `{"message":"a problem appeared"}`)
		}
	} else {
		utils.SendResponse(res, http.StatusForbidden, `{"message":"user was not connected"}`)
	}
}

func AddUser(res http.ResponseWriter, req *http.Request) {
	mail := req.FormValue("mail")
	login := req.FormValue("login")
	password := req.FormValue("password")
	if existingLogin(login) {
		utils.SendResponse(res, http.StatusForbidden, `{"message":"problem login alredy exist"}`)
	} else {
		if insertUser(login, password, mail) {
			utils.SendResponse(res, http.StatusOK, `{"message":"New user created"}`)
		} else {
			utils.SendResponse(res, http.StatusInternalServerError, `{"message":"a problem appeared"}`)
		}
	}
}

func ModifyUser(_ http.ResponseWriter, _ *http.Request) {

}

func DeleteUser(_ http.ResponseWriter, _ *http.Request) {

}

func IsUser(login string, password string) bool {
	return true
}

func existingLogin(login string) bool {
	return false
}

func searchUser(login string) User {
	db := database.Connect()
	if db == nil {
		return nil
	}
	u := User{}
	db.QueryRow("Select login, mail, cagnotte From User where login=?;", login).Scan(&.login, &u.mail, &u.cagnotte)
	db.Close()
	return u
}

func insertUser(login string, password string, mail string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	
	res, err := db.Exec("INSERT INTO Utilisateur VALUES (?, ?, ?, 100);", login, password, mail)
	db.Close()

	r, _ := res.RowsAffected()
	if r == 0 || err != nil {
		return false
	} else {
		return true
	}
}
