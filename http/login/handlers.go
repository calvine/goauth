package login

import "net/http"

func GetLoginHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Sup"))
}

func PostLoginHandler(rw http.ResponseWriter, r *http.Request) {

}
