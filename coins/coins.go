package coins

import (
	"gitlab.com/CTeillet/pc3r-projet/connexion"
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"net/http"
)

func Generate(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")
	montant	  := r.FormValue("montant")

	if login:=connexion.IsConnected(idSession)!="" {
		db := database.Connect()
		if db == nil {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection database"}`)
		}

		res, err := db.Exec("Update Utilisateur SET cagnotte = cagnotte + 100 where login=? ;", login)
		if res.RowsAffected()==1 && err=nil {

		}
		
		

	}
}

func addMoney(login bool, montant string) interface{} {

}
