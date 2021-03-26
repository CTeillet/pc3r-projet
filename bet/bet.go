package bet

import "net/http"

type Bet struct {
	id             int
	idMatch        int
	equipeGagnante string
	cote           float32
	montant        int
}

func GetBet(_ http.ResponseWriter, _ *http.Request) {

}

func AddBet(_ http.ResponseWriter, _ *http.Request) {

}

func DeleteBet(_ http.ResponseWriter, _ *http.Request) {

}
