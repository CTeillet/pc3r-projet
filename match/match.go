package match

import (
	"database/sql"
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
	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"`)
		return
	}

	db := database.Connect()
	if db == nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
		return
	}

	var res *sql.Rows
	var err error
	if req == "" {
		res, err = db.Query("Select * From Match where status='open';")
	} else {
		res, err = db.Query("Select * From Match where status='open' and (sport=? or region=? or equipeA=? or equipeB=?);", req)
	}

	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
		return
	}
	err = db.Close()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"`)
		return
	}
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
	utils.SendResponse(w, http.StatusOK, `{"message": "result of match", "result":`+string(resultJSON)+"}")

}

//Ne pas appeler : LoadAllPastMatch
func _() {
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
		go JSONMatch2SQL(resp)
	}
}

func LoadComingMatchFor2Week() {
	req := "https://api.pandascore.co/lol/matches/upcoming?token=4xg85-0CNl9sOdk-tyFooufCsE8qchuK478B5bUoAOV0j3cREdQ"
	t := time.Now()
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
		max = 0
	}

	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		fmt.Println(s)
		resp, _ := http.Get(s)
		go JSONMatch2SQL(resp)
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
			addMatch(v.Videogame.Name, v.League.Name, "", "", v.Status, "", v.BeginAt)
		}
	}
}

func addMatch(sport string, league string, equipeA string, equipeB string, statut string, winner string, date time.Time) {
	db := database.Connect()
	//_, err := db.Exec("Insert into `Match` (sport, league, equipeA, equipeB, cote,statut, vainqueur, date) VALUES (?, ?, ?, ?, 1.0, ?, ?, ?);", sport, league, equipeA, equipeB, statut, winner, date)
	//fmt.Printf("Update `Match` set equipeA=%v and equipeB=%v and vainqueur=%v and statut=%v where sport=%v and league=%v and equipeA='' and equipeB='' and date=%v ;\n", equipeA, equipeB, winner, statut, sport, league, date)
	_, err := db.Exec("Update `Match` set equipeA=? and equipeB=? and vainqueur=? and statut=? where sport=? and league=? and equipeA='' and equipeB='' and date=? ;", equipeA, equipeB, winner, statut, sport, league, date)
	if err != nil {
		_, err := db.Exec("Insert into `Match` (sport, league, equipeA, equipeB, cote,statut, vainqueur, date) VALUES (?, ?, ?, ?, 1.0, ?, ?, ?);", sport, league, equipeA, equipeB, statut, winner, date)
		if err != nil {
			if !strings.Contains(err.Error(), "Duplicate") {
				panic(err.Error())
			}

		}
	}
	err = db.Close()
}

func LoadResultMatchFor1Hour() {
	req := "https://api.pandascore.co/lol/matches/past?token=4xg85-0CNl9sOdk-tyFooufCsE8qchuK478B5bUoAOV0j3cREdQ"
	t := time.Now()
	req += "&range[begin_at]=" + strings.Split(t.Add(-1*time.Hour).Format("2006-01-02T15:04:05-0700"), "+")[0] + "," + strings.Split(t.Format("2006-01-02T15:04:05-0700"), "+")[0]
	s := req + "&page[size]=100"
	fmt.Println(s)
	resp, _ := http.Get(s)
	JSONMatchUpdate(resp)

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
		max = 0
	}
	fmt.Println(max)
	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		fmt.Println(s)
		resp, _ := http.Get(s)
		JSONMatchUpdate(resp)
	}

}

func JSONMatchUpdate(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	var data utils.MatchJSON // TopTracks
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	updateMulipleMatch(data)
}

func updateMulipleMatch(data utils.MatchJSON) {
	for _, v := range data {
		if len(v.Opponents) == 2 {
			updateMatch(v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Acronym, v.Opponents[1].Opponent.Acronym, v.Winner.Acronym, v.ScheduledAt, v.Status)
		}
	}
}

func updateMatch(sport string, league string, equipeA string, equipeB string, winner string, date time.Time, statut string) {
	db := database.Connect()
	fmt.Printf("Update `projet-pc3r`.`Match` SET `vainqueur`=%v and `statut`=%v where sport=%v and league=%v and equipeA=%v and equipeB=%v and `date`=%v and statut='not_started';\n", winner, statut, sport, league, equipeA, equipeB, date)
	_, err := db.Exec("Update `projet-pc3r`.`Match` SET `vainqueur`=? , `statut`=? where sport=? and league=? and equipeA=? and equipeB=? and `date`=? and statut='not_started';", winner, statut, sport, league, equipeA, equipeB, date)
	if err != nil {
		panic(err.Error())
	}
	err = db.Close()
}
