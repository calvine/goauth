package handlers

import "net/http"

func (hh HttpHandler) GetLoginHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Sup"))
}

func (hh HttpHandler) PostLoginHandler(rw http.ResponseWriter, r *http.Request) {

}
