package handlers

import (
	"gopkg.in/yaml.v3"
	"net/http"
)

// handlerNotFound, хендлер начальной страницы сервера
func (srv *Server) handleFunc(rw http.ResponseWriter, rq *http.Request) {

	out, _ := yaml.Marshal(srv.Setting)

	if _, err := rw.Write(out); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// handlerNotFound, хендлер начальной страницы сервера
func (srv *Server) handleFuncUpdateDB(rw http.ResponseWriter, rq *http.Request) {
	srv.InitDB()

	out, _ := yaml.Marshal(srv.Setting)

	if _, err := rw.Write(out); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
