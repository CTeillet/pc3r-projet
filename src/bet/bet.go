package bet

import (
	"database/sql"
	"encoding/json"
	"errors"
	"gitlab.com/CTeillet/pc3r-projet/src/database"
	"gitlab.com/CTeillet/pc3r-projet/src/match"
	"gitlab.com/CTeillet/pc3r-projet/src/user"
	"gitlab.com/CTeillet/pc3r-projet/src/utils"
	"net/http"
	"strconv"
	"time"
)

type Bet struct {
	id             int
	idMatch        int
	equipeGagnante string
	cote           float32
	montant        float32
	login          string
	resultat       string
	date           time.Time
}

func GetBet(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")

	login := utils.IsConnected(idSession)
	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"`)
		return
	}

	db := database.Connect()
	if db == nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
		return
	}

	res, err := db.Query("Select * From `projet-pc3r`.`Pari` where login=?;", login)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
		return
	}
	err = db.Close()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
		return
	}
	resultat := make([]Bet, 0)
	for res.Next() {
		b := Bet{}
		err := res.Scan(&b.id, &b.equipeGagnante, &b.cote, &b.montant, &b.login, &b.resultat, &b.date)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem reading result request"`)
			return
		}
		resultat = append(resultat, b)
	}
	resultJSON, err := json.Marshal(resultat)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem creation of JSON"`)
		return
	}
	utils.SendResponse(w, http.StatusOK, `{"message": "request effected", "result":`+string(resultJSON)+"}")
}

func AddBet(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")
	idMatch := r.FormValue("idMatch")
	equipeGagnante := r.FormValue("equipeGagnante")
	cote, err := strconv.ParseFloat(r.FormValue("cote"), 32)
	if err != nil {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "wrong value for cote"`)
	}
	montant, err := strconv.ParseFloat(r.FormValue("montant"), 32)
	if err != nil {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "wrong value for montant"`)
	}

	login := utils.IsConnected(idSession)

	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"`)
		return
	}

	testInsert := addBetSql(idMatch, equipeGagnante, float32(cote), float32(montant), login)

	if !testInsert {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
	} else {
		user.AlterMoney(login, float32(montant))
		utils.SendResponse(w, http.StatusOK, `{"message":"New bet created"}`)
	}

}

func DeleteBet(w http.ResponseWriter, r *http.Request) {
	idPari := r.FormValue("idPari")
	idSession := r.FormValue("idSession")
	login := utils.IsConnected(idSession)

	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"`)
		return
	}

	testInsert := removeBetSQL(idPari, login)

	if !testInsert {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
	} else {
		utils.SendResponse(w, http.StatusOK, `{"message":"New user created"}`)
	}
}

func removeBetSQL(pari, login string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	res, err := db.Exec("Delete from `projet-pc3r`.Pari where `projet-pc3r`.Pari.id=? and `projet-pc3r`.Pari.login=?;", pari, login)
	if err != nil {
		return false
	}
	test, err := res.RowsAffected()
	if err != nil || test != 1 {
		return false
	}
	return true
}

func addBetSql(idMatch, equipeGagnante string, cote float32, montant float32, login string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	res, err := db.Exec("Insert into `projet-pc3r`.Pari(`projet-pc3r`.Pari.idmatch, `projet-pc3r`.Pari.equipegagnante, `projet-pc3r`.Pari.cote, `projet-pc3r`.Pari.montant, `projet-pc3r`.Pari.login) Values(?, ?, ?, ?, ?) ;", idMatch, equipeGagnante, cote, montant, login)
	if err != nil {
		return false
	}
	test, err := res.RowsAffected()
	if err != nil || test != 1 {
		return false
	}
	return true
}

func UpdateResult1Hour() {
	db := database.Connect()
	if db == nil {
		panic(errors.New("problem database connection"))
	}
	res, err := db.Query("Select * from  Pari as P where resultat='coming' and EXISTS( Select * From `Match` where id=P.idMatch and statut='finished');")
	if err != nil {
		panic(err.Error())
	}
	for res.Next() {
		b := Bet{}
		err := res.Scan(&b.id, &b.idMatch, &b.equipeGagnante, &b.cote, &b.montant, &b.login, &b.resultat, &b.date)

		if err != nil {
			return
		}

		var r sql.Result
		win := match.WinnerIdMatch(b.idMatch)
		if b.equipeGagnante == win {
			user.AlterMoney(b.login, b.montant*b.cote)
			r, err = db.Exec("Update Pari set resultat='win' where id=?", b.id)
		} else {
			r, err = db.Exec("Update Pari set resultat='loose' where id=?", b.id)
		}
		if err != nil {
			return
		}
		res, err := r.RowsAffected()
		if res != 1 {
			return
		}
	}

}
