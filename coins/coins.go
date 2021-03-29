package coins

import (
	"gitlab.com/CTeillet/pc3r-projet/connexion"
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"net/http"
)

func Generate(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")
	montant := r.FormValue("montant")

	login := connexion.IsConnected(idSession)
	if login != "" {
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

		r, err := res.RowsAffected()
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection database"}`)
			return
		}

		if r != 1 {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with request"}`)
			return
		}
		utils.SendResponse(w, http.StatusOK, `{"message": "user has now more coins"}`)
		return

	} else {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"}`)
		return
	}
}
