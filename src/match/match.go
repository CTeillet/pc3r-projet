package match

import (
	"database/sql"
	"encoding/json"
	"gitlab.com/CTeillet/pc3r-projet/src/database"
	"gitlab.com/CTeillet/pc3r-projet/src/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Match struct {
	Id        int       `json:"id"`
	Sport     string    `json:"sport"`
	League    string    `json:"league"`
	EquipeA   string    `json:"equipeA"`
	EquipeB   string    `json:"equipeB"`
	Cote      float32   `json:"cote"`
	Statut    string    `json:"statut"`
	Vainqueur string    `json:"vainqueur"`
	Date      time.Time `json:"date"`
}

func GetMatch(w http.ResponseWriter, r *http.Request) {
	req := r.FormValue("req")
	idSession := r.FormValue("idSession")

	login := utils.IsConnectedIdSession(idSession)
	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"}`)
		return
	}

	db := database.Connect()
	if db == nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection to database"}`)
		return
	}

	var res *sql.Rows
	var err error
	if req == "" {
		res, err = db.Query("Select * From `Match` where statut='not_started' and equipeA<>'' and equipeB<>'' order by date;")
	} else {
		res, err = db.Query("Select * From `Match` where statut='not_started' and (sport=? or league=? or equipeA=? or equipeB=?) order by date;", req, req, req, req)
	}

	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with searching database"}`)
		return
	}
	err = db.Close()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with closing database"}`)
		return
	}
	resultat := make([]Match, 0)
	for res.Next() {
		m := Match{}
		err := res.Scan(&m.Id, &m.Sport, &m.League, &m.EquipeA, &m.EquipeB, &m.Cote, &m.Statut, &m.Vainqueur, &m.Date)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem reading result request"}`)
			return
		}
		resultat = append(resultat, m)
	}
	resultJSON, err := json.Marshal(resultat)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem creation of JSON"}`)
		return
	}
	utils.SendResponse(w, http.StatusOK, `{"message": "coming matches", "result":`+string(resultJSON)+"}")

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
	//fmt.Println(max)
	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		//fmt.Println(s)
		resp, _ := http.Get(s)
		go JSONMatch2SQL(resp)
	}
}

func LoadComingMatchFor2Week() {
	req := "https://api.pandascore.co/lol/matches/upcoming?token=4xg85-0CNl9sOdk-tyFooufCsE8qchuK478B5bUoAOV0j3cREdQ"
	t := time.Now()
	req += "&range[begin_at]=" + strings.Split(t.Format("2006-01-02T15:04:05-0700"), "+")[0] + "," + strings.Split(t.Add(time.Hour*24*7*2).Format("2006-01-02T15:04:05-0700"), "+")[0]
	s := req + "&page[size]=100"
	//fmt.Println(s)
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
	//fmt.Println(max)
	if err != nil {
		max = 0
	}

	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		//fmt.Println(s)
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
		//fmt.Println(v)
		if len(v.Opponents) == 2 {
			addMatch(v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Acronym, v.Opponents[1].Opponent.Acronym, v.Status, v.Winner.Acronym, v.OriginalScheduledAt)
		} else {
			addMatch(v.Videogame.Name, v.League.Name, "", "", v.Status, "", v.OriginalScheduledAt)
		}
	}
}

func addMatch(sport string, league string, equipeA string, equipeB string, statut string, winner string, date time.Time) {
	db := database.Connect()
	//_, err := db.Exec("Insert into `Match` (sport, league, equipeA, equipeB, cote,statut, vainqueur, date) VALUES (?, ?, ?, ?, 1.0, ?, ?, ?);", sport, league, equipeA, equipeB, statut, winner, date)
	//fmt.Printf("Update `Match` set equipeA=%v , equipeB=%v , vainqueur=%v , statut=%v where sport=%v and league=%v and equipeA='' and equipeB='' and date=%v ;\n", equipeA, equipeB, winner, statut, sport, league, date)
	r, err := db.Exec("Update `Match` set equipeA=? , equipeB=? , vainqueur=? , statut=? where sport=? and league=? and equipeA='' and equipeB='' and date=? ;", equipeA, equipeB, winner, statut, sport, league, date)
	if err == nil {
		nbRows, err2 := r.RowsAffected()
		if err2 != nil || nbRows != 1 {
			//fmt.Println(err.Error())
			_, err := db.Exec("Insert into `Match` (sport, league, equipeA, equipeB, cote,statut, vainqueur, date) VALUES (?, ?, ?, ?, 1.0, ?, ?, ?);", sport, league, equipeA, equipeB, statut, winner, date)
			if err != nil {
				if !strings.Contains(err.Error(), "Duplicate") {
					panic(err.Error())
				}

			}
		}
	}
	err = db.Close()
}

func LoadResultMatchFor3Hours() {
	req := "https://api.pandascore.co/lol/matches/past?token=4xg85-0CNl9sOdk-tyFooufCsE8qchuK478B5bUoAOV0j3cREdQ"
	t := time.Now()
	req += "&range[end_at]=" + strings.Split(t.Add(-3*time.Hour).Format("2006-01-02T15:04:05-0700"), "+")[0] + "," + strings.Split(t.Format("2006-01-02T15:04:05-0700"), "+")[0]
	s := req + "&page[size]=100"
	//fmt.Println(s)
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
	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		//fmt.Println(s)
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
			//fmt.Printf("Sport :%v , League :%v ,  OppentsA : %v, OpponentsB : %v, vainqueur : %v, date : %v, statut : %v\n", v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Acronym, v.Opponents[1].Opponent.Acronym, v.Winner.Acronym, v.ScheduledAt, v.Status)
			updateMatch(v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Acronym, v.Opponents[1].Opponent.Acronym, v.Winner.Acronym, v.OriginalScheduledAt, v.Status)
		}
	}
}

func updateMatch(sport string, league string, equipeA string, equipeB string, winner string, date time.Time, statut string) {
	db := database.Connect()

	_, err := db.Exec("Update `projet-pc3r`.`Match` SET `vainqueur`=? , `statut`=? where sport=? and league=? and equipeA=? and equipeB=? and `date`=? and statut='not_started';", winner, statut, sport, league, equipeA, equipeB, date)
	if err != nil {
		panic(err.Error())
	}
	err = db.Close()
}

func WinnerIdMatch(idMatch int) string {
	db := database.Connect()
	if db == nil {
		return ""
	}
	m := Match{}
	err := db.QueryRow("Select * From `Match` where id=?;", idMatch).Scan(&m.Id, &m.Sport, &m.League, &m.EquipeA, &m.EquipeB, &m.Cote, &m.Statut, &m.Vainqueur, &m.Date)
	if err != nil {
		panic(err.Error())
	}
	return m.Vainqueur

}
