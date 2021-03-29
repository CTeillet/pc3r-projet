package user

import (
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"strconv"
	"net/http"
)

type User struct {
	login    string
	mail     string
	cagnotte int
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	//idSession := req.FormValue("idSession")
	login := req.FormValue("login")
	if true{//connexion.IsConnected(idSession) {
		var user *User = nil
		searchUser(login, user)
		if user != nil {
			if user.login != "" {
				cagnotte := strconv.Itoa(user.cagnotte)
				utils.SendResponse(res, http.StatusOK, `{"message":"user found", "login":"`+user.login+`", "mail":"`+user.mail+`", "cagnotte":"`+cagnotte+`"}`)
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

func searchUser(login string, u *User)  {
	db := database.Connect()
	if db == nil {
		return
	}
	
	db.QueryRow("Select login, mail, cagnotte From User where login=?;", login).Scan(u.login, u.mail, u.cagnotte)
	db.Close()
	return
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
