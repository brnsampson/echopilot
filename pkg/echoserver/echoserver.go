package echoserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Section 1: application specific functionas and behavior. This could potentially be moved into
// a separate file if it gets extensive enough.s
func Echo(s string) (string, error) {
	return s, nil
}

func (s *server) EchoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.FormValue("data")
		resp, err := Echo(data)
		if err != nil {
			s.logger.Errorf("echoHandler: %v", err)
			http.NotFound(w, r)
		} else {
			fmt.Fprintf(w, resp)
		}
		return
	}
}

// This interface can and will probably need to be modified if you don't want to use the uber zap logger.
type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Sync() error
}

// Section 2: common server implementation boilerplate
type server struct {
	logger   Logger
	router   *http.ServeMux
	wg       sync.WaitGroup
	err      chan error
	done     chan struct{}
	stop     chan os.Signal
	sig      chan os.Signal
	exitCode int
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ServeWithReload() {
	for {
		addr := os.Getenv("ECHO_ADDR")
		if addr == "" {
			addr = "127.0.0.1:8080"
		}
		serv := &http.Server{Addr: addr, Handler: s.router}

		go func(err chan<- error) {
			s.logger.Infof("Listening on %s", addr)
			if e := serv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
				err <- e
			}
			return
		}(s.err)

		select {
		case <-s.sig:
			s.logger.Info("SIGHUP recieved. Reloading...")
			begin := time.Now()
			s.logger.Debug("Halting Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := serv.Shutdown(ctx); err != nil {
				s.logger.Debugf("Failed to gracefully shutdown server: %v", err)
			} else {
				s.logger.Debugf("Server halted in %v", time.Since(begin))
			}
		case <-s.done:
			s.logger.Info("Server shutting down...")
			begin := time.Now()
			s.logger.Debug("Halting Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := serv.Shutdown(ctx); err != nil {
				s.logger.Debugf("Failed to gracefully shutdown server: %v", err)
			} else {
				s.logger.Debugf("Server halted in %v", time.Since(begin))
			}
			s.wg.Done()
			return
		}
	}
}

func (s *server) Run() {
	s.wg.Add(1)
	go s.waitOnInterrupt()
	go s.ServeWithReload()
	return
}

func (s *server) BlockingRun() int {
	s.wg.Add(1)
	go s.ServeWithReload()
	s.waitOnInterrupt()
	return s.exitCode
}

func (s *server) waitOnInterrupt() {
	select {
	case <-s.stop:
		s.logger.Info("Interrupt/kill recieved. Exiting...")
		s.exitCode = 0
		s.Shutdown()
	case e := <-s.err:
		s.logger.Errorf("Encoutered error %v. Exiting...", e)
		s.Shutdown()
		s.exitCode = 1
	}
}

func (s *server) Shutdown() {
	signal.Stop(s.sig)
	signal.Stop(s.stop)
	close(s.done)

	s.wg.Wait()
	s.logger.Info("All waits done. Server execution complete.")

	close(s.stop)
	close(s.sig)
	close(s.err)
	return
}

func NewServer(logger Logger) *server {
	// First set up signal hanlding so that we can reload and stop.
	hups := make(chan os.Signal, 1)
	stop := make(chan os.Signal, 1)

	signal.Notify(hups, syscall.SIGHUP)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var waitgroup sync.WaitGroup

	done := make(chan struct{}, 1)
	err := make(chan error, 1)

	m := http.NewServeMux()

	// If exitCode is -1 then execution has not completed yet.
	server := &server{logger, m, waitgroup, err, done, stop, hups, -1}
	server.routes()
	return server
}
