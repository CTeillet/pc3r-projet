package match

import (
	"encoding/json"
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"net/http"
)

type Match struct {
	id      int
	sport   string
	region  string
	equipeA string
	equipeB string
	cote    float32
	statut  string
}

func GetMatch(w http.ResponseWriter, r *http.Request) {
	req := r.FormValue("login")
	idSession := r.FormValue("idSession")

	login := utils.IsConnected(idSession)
	if login != "" {
		if req == "" {
			db := database.Connect()
			if db != nil {
				res, err := db.Query("Select * From Match where status='open';")
				if err != nil {
					utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
					return
				}
				err = db.Close()
				if err != nil {
					utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
				} else {
					resultat := []Match{}
					for res.Next() {
						m := Match{}
						err := res.Scan(&m.id, &m.sport, &m.region, &m.equipeA, &m.equipeB, &m.cote, &m.statut)
						if err != nil {
							utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem reading result request"`)
							return
						}
						resultat = append(resultat, m)
					}
					resultJSON, err := json.Marshal(resultat)
					if err != nil {
						utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem creation of JSON"`)
						return
					}
					utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem creation of JSON", "resultat":`+string(resultJSON)+"}")
				}
			} else {
				utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
			}
		} else {
			db := database.Connect()
			if db != nil {
				res, err := db.Query("Select * From Match where status='open' and (sport=? or region=? or equipeA=? or equipeB=?);", req)
				if err != nil {
					utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
					return
				}
				err = db.Close()
				if err != nil {
					utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
				} else {
					resultat := []Match{}
					for res.Next() {
						m := Match{}
						err := res.Scan(&m.id, &m.sport, &m.region, &m.equipeA, &m.equipeB, &m.cote, &m.statut)
						if err != nil {
							utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem reading result request"`)
							return
						}
						resultat = append(resultat, m)
					}
					resultJSON, err := json.Marshal(resultat)
					if err != nil {
						utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem creation of JSON"`)
						return
					}
					utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem creation of JSON", "resultat":`+string(resultJSON)+"}")
				}
			} else {
				utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
			}
		}
	} else {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"`)
	}
}
