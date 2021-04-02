package user

import (
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"net/http"
	"strconv"
)

type User struct {
	login    string
	mail     string
	cagnotte int
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	//Recuperation des parametres de la requete http
	idSession := req.FormValue("idSession")
	login := req.FormValue("login")
	//verif connexion 
	if utils.IsConnected(idSession) != "" {
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
	// TODO
}

func addDailyCoins(_ http.ResponseWriter, _ *http.Request) {

}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	login := req.FormValue("login")
	password := req.FormValue("password")
	idSession :=req.FormValue("idSession")
	if utils.IsConnected(idSession) != "" && idSession != "" {
		if utils.isUser(login, password) {
			if removeUser(login, password) {
				utils.SendResponse(res, http.StatusOK, `{"message":"deleted user"}`)
			} else {
				utils.SendResponse(res, http.StatusInternalServerError, `{"message":"problem with database"}`)
			}

		} else {
			utils.SendResponse(res, http.StatusForbidden, `{"message":"Error : wrong login or password"}`)
		}
	}
}

func existingLogin(login string) bool {
	// TODO
	return false
}

func searchUser(login string, u *User) {
	db := database.Connect()
	if db == nil {
		return
	}

	db.QueryRow("Select login, mail, cagnotte From Utilisateur where login=?;", login).Scan(u.login, u.mail, u.cagnotte)
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

func removeUser(login string, password string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	res, err := db.Exec("Delete from Utilisateur where login=? and password=?;", login, password)
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}

	r, _ := res.RowsAffected()
	if r == 1 {
		return true
	}
	return false
}

func isUser(login string, password string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	count := 0
	err := db.QueryRow("Select count(*) From Connexion where login=? and password=?;", login, password).Scan(&count)
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}
	if count == 1{
		return true
	}
	return false
}