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
		var user = searchUser(login)
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
	if acceptLogin(login) {
		if insertUser(login, password, mail) {
			utils.SendResponse(res, http.StatusOK, `{"message":"New user created"}`)
		} else {
			utils.SendResponse(res, http.StatusInternalServerError, `{"message":"a problem appeared"}`)
		}
	} else {
		utils.SendResponse(res, http.StatusForbidden, `{"message":"problem login alredy exist"}`)
	}
}

func ModifyUser(_ http.ResponseWriter, _ *http.Request) {
	// TODO
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	login := req.FormValue("login")
	password := req.FormValue("password")
	idSession := req.FormValue("idSession")
	if utils.IsConnected(idSession) != "" && idSession != "" {
		if utils.IsUser(login, password) {
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

func acceptLogin(login string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	var count int
	err := db.QueryRow("Select count(*) From Utilisateur where login=?;", login).Scan(&count)
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}
	if count == 0 {
		return true
	}
	return false
}

func searchUser(login string) *User {
	db := database.Connect()
	if db == nil {
		return nil
	}

	u := User{}

	db.QueryRow("Select login, mail, cagnotte From Utilisateur where login=?;", login).Scan(&u.login, &u.mail, &u.cagnotte)
	db.Close()
	return &u
}

func insertUser(login string, password string, mail string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}

	res, err := db.Exec("INSERT INTO Utilisateur VALUES (?, ?, ?, 100);", login, password, mail)
	err = db.Close()
	if err != nil {
		return false
	}

	r, err := res.RowsAffected()
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
