package match

import (
	"encoding/json"
	"fmt"
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"io/ioutil"
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

func UpdateMatch() {
	//t := time.Now()
	s := ""
	//for i:=0; i<1 ; i++{
	//	s+=t.Format("YYYY-MM-DD")
	//	t = t.Add(24*time.Hour)
	//	if i!= 0 {
	//		s+=",%20"
	//	}
	//}
	//s = t.Format("2006-01-02")
	s = "2021-03-27"
	fmt.Println(s)
	fmt.Println("REQUESTEEE    " + "https://fly.sportsdata.io/v3/lol/scores/json/GamesByDate/" + s + "?key=c86e4989da6247358a15b0c3ab5dbe66")
	resp, _ := http.Get("https://fly.sportsdata.io/v3/lol/scores/json/GamesByDate/" + s + "?key=c86e4989da6247358a15b0c3ab5dbe66")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	var data interface{} // TopTracks
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Results: %v\n", data)
}
