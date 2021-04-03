package match

import (
	"encoding/json"
	"fmt"
	"gitlab.com/CTeillet/pc3r-projet/database"
	"gitlab.com/CTeillet/pc3r-projet/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
					resultat := make([]Match, 0)
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
					resultat := make([]Match, 0)
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

//Ne pas appeler : LoadAllPastMatch
func LoadAllPastMatch() {
	req := "https://api.pandascore.co/lol/matches/past?token=4xg85-0CNl9sOdk-tyFooufCsE8qchuK478B5bUoAOV0j3cREdQ"

	resp, _ := http.Get(req + "&page[size]=100")
	JSONMatch2SQL(resp)

	test := resp.Header.Get("Link")
	res := strings.Split(test, ",")
	last := ""
	for _, v := range res {
		if strings.Contains(v, "last") {
			last = strings.Split(v, ";")[0][2 : len(strings.Split(v, ";")[0])-1]
		}
	}

	u, err := url.Parse(last)
	if err != nil {
		panic(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		panic(err)
	}
	fmt.Println(max)
	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		fmt.Println(s)
		resp, _ := http.Get(s)
		JSONMatch2SQL(resp)
		time.Sleep(100 * time.Millisecond)
	}
}

func LoadComingMatchFor2Week() {
	req := "https://api.pandascore.co/lol/matches/upcoming?token=4xg85-0CNl9sOdk-tyFooufCsE8qchuK478B5bUoAOV0j3cREdQ"
	t := time.Now().Add(time.Minute)
	req += "&range[begin_at]=" + strings.Split(t.Format("2006-01-02T15:04:05-0700"), "+")[0] + "," + strings.Split(t.Add(time.Hour*24*7*2).Format("2006-01-02T15:04:05-0700"), "+")[0]
	s := req + "&page[size]=100"
	fmt.Println(s)
	resp, _ := http.Get(s)
	JSONMatch2SQL(resp)

	test := resp.Header.Get("Link")
	res := strings.Split(test, ",")
	last := ""
	for _, v := range res {
		if strings.Contains(v, "last") {
			last = strings.Split(v, ";")[0][2 : len(strings.Split(v, ";")[0])-1]
		}
	}

	u, err := url.Parse(last)
	if err != nil {
		panic(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(q.Get("page"))
	fmt.Println(max)
	if err != nil {
		//panic(err.Error())
		max = 0
	}

	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		fmt.Println(s)
		resp, _ := http.Get(s)
		go JSONMatch2SQL(resp)
		time.Sleep(100 * time.Millisecond)
	}
}

func JSONMatch2SQL(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	var data utils.MatchJSON // TopTracks
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	addMulipleMatch(data)
}

func addMulipleMatch(data utils.MatchJSON) {
	for _, v := range data {
		//time.Sleep(150*time.Millisecond)
		if len(v.Opponents) == 2 {
			addMatch(v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Acronym, v.Opponents[1].Opponent.Acronym, v.Status, v.Winner.Acronym, v.BeginAt)
		} else {
			addMatch(v.Videogame.Name, v.League.Name, "", "", v.Status, v.Winner.Acronym, v.BeginAt)
		}
	}
}

func addMatch(sport string, league string, equipeA string, equipeB string, statut string, winner string, date time.Time) {
	db := database.Connect()
	_, err := db.Exec("Insert into `Match` (sport, league, equipeA, equipeB, cote,statut, vainqueur, date) VALUES (?, ?, ?, ?, 1.0, ?, ?, ?);", sport, league, equipeA, equipeB, statut, winner, date)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		panic(err.Error())
	}
	db.Close()
}
