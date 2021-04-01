package bet

import (
	"net/http"
	"time"
)

type Bet struct {
	id             int
	idMatch        int
	equipeGagnante string
	cote           float32
	montant        int
	date           time.Time
}

func GetBet(_ http.ResponseWriter, _ *http.Request) {

}

func AddBet(_ http.ResponseWriter, _ *http.Request) {

}

func DeleteBet(_ http.ResponseWriter, _ *http.Request) {

}
