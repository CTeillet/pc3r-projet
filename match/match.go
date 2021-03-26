package match

import "net/http"

type Match struct {
	id      int
	equipeA string
	equipeB string
	cote    float32
}

func GetMatch(_ http.ResponseWriter, _ *http.Request) {

}
