package match

import (
	"encoding/json"
	"fmt"
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"io/ioutil"
	"net/http"
	"time"
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

func UpdateMatchPast() {

	s := "https://api.pandascore.co/lol/matches/past?page[size]=100&token=4xg85-0CNl9sOdk-tyFooufCsE8qchuK478B5bUoAOV0j3cREdQ"
	resp, _ := http.Get(s)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	var data utils.MatchPastJSON // TopTracks
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(len(data))
	fmt.Println(resp.Header)
	for _, v := range data {
		addMatch(v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Acronym, v.Opponents[1].Opponent.Acronym, v.Winner.Acronym, v.BeginAt)
	}

}

func addMatch(sport string, league string, equipeA string, equipeB string, winner string, date time.Time) {
	db := database.Connect()
	_, err := db.Exec("Insert into `Match` (sport, league, equipeA, equipeB, cote,statut, vainqueur, date) VALUES (?, ?, ?, ?, 1.0, 'open', ?, ?);", sport, league, equipeA, equipeB, winner, date)
	if err != nil {
		panic(err.Error())
	}
	db.Close()
}
