package coins

import (
	"gitlab.com/CTeillet/pc3r-projet/src/database"
	"gitlab.com/CTeillet/pc3r-projet/src/utils"
	"net/http"
)

func Generate(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")
	montant := r.FormValue("montant")

	login := utils.IsConnected(idSession)
	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"}`)
		return
	}

	db := database.Connect()
	if db == nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection database"}`)
		return
	}

	res, err := db.Exec("Update Utilisateur SET cagnotte = cagnotte + ? where login=? ;", montant, login)

	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection database"}`)
		return
	}

	row, err := res.RowsAffected()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection database"}`)
		return
	}

	if row != 1 {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with request"}`)
		return
	}
	utils.SendResponse(w, http.StatusOK, `{"message": "user has now more coins"}`)

}
