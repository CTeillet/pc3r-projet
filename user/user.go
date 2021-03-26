package user

import "net/http"

type User struct {
	id       int
	login    string
	password string
	mail     string
}

func GetUser(_ http.ResponseWriter, _ *http.Request) {

}

func AddUser(_ http.ResponseWriter, _ *http.Request) {

}

func ModifyUser(_ http.ResponseWriter, _ *http.Request) {

}

func DeleteUser(_ http.ResponseWriter, _ *http.Request) {

}

func IsUser(login string, password string) bool {
	return true
}
