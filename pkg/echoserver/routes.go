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

func IndexHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			filePath := "dist/" + r.URL.Path
			http.ServeFile(w, r, filePath)
		} else {
			http.ServeFile(w, r, entrypoint)
		}
	}
}

func (s *server) routes() {
	fs := http.FileServer(http.Dir("dist"))
	s.router.HandleFunc("/api/echo", WrapHandler(s.EchoHandler(), s.logger))
	s.router.Handle("/dist/", http.StripPrefix("/dist/", fs))
	s.router.HandleFunc("/", WrapHandler(IndexHandler("dist/index.html"), s.logger))
}
