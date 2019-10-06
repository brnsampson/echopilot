package echoserver

import (
	"net/http"
	"time"
)

func WrapHandler(fn http.HandlerFunc, logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		fn(w, r)
		logger.Debugf("Served request for %s in %v", r.URL.Path, time.Since(begin))
		logger.Debugf("Body data field: %s", r.FormValue("data"))
		return
	}
}

func (s *server) routes() {
	s.router.HandleFunc("/", WrapHandler(s.EchoHandler(), s.logger))
}
