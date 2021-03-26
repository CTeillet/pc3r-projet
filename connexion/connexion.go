package connexion

import (
	"gitlab.com/CTeillet/pc3r-projet/user"
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
		idSession := addConnexion(login, password)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"unknown user", "idSession":"` + idSession + `"}`))
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"unknown user"}`))
	}
}

func Disconnect(_ http.ResponseWriter, r *http.Request) {

}

func IsConnected(idSession string) bool {
	//TODO : interaction avec la base de données : verifier si l'utilisateur est connecté,
	// s'il est connecté vérifier de quand date la dernière interaction si elle supérieur à X minutes
	// deconecté l'utilisateur et renvoyé false, sinon mettre à jour la date et renvoyé true
	// sinon renvoyé false
	return false
}

func addConnexion(login string, password string) string {
	//TODO : ajouter la connexion à la base de données et renvoyé l'idSession génére
	return "test"
}
