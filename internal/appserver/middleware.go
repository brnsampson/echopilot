package appserver

import (
	"net/http"
	"time"
    "github.com/charmbracelet/log"
)

// Helper structs for pilfering information from the wrapped handlers

type responseSpy struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func (s responseSpy) Write(buf []byte) (int, error) {
	written, e := s.ResponseWriter.Write(buf)
	if e != nil {
		return written, e
	}
	s.bytesWritten += int64(written)
	return written, e
}

func (s responseSpy) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

// Middleware struct for adding logging to requests
func NewLoggingHandler(toWrap http.Handler, logger *log.Logger) *loggingHandler {
	return &loggingHandler{toWrap, logger}
}

type loggingHandler struct {
	wrappedHandler http.Handler
	logger         *log.Logger
}

func (l *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	spy := responseSpy{w, 0, 0}
	begin := time.Now()
	l.wrappedHandler.ServeHTTP(spy, r)
	l.logger.Infof("Status Code %d for %s at %s in %+v", spy.statusCode, r.Method, r.URL.Path, time.Since(begin))
	l.logger.Debugf("Replying to %s request with: %s", r.URL.Path)
}
