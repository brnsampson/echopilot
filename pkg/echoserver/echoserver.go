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

type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Sync() error
}

type comms struct {
	wg   *sync.WaitGroup
	err  chan<- error
	done <-chan struct{}
	sig  <-chan os.Signal
}

func Echo(s string) (string, error) {
	return s, nil
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("data")
	resp, err := Echo(data)
	if err != nil {
		fmt.Fprint(w, "404")
	} else {
		fmt.Fprintf(w, resp)
	}
}

func MakeHandler(fn func(http.ResponseWriter, *http.Request), logger Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		fn(w, r)
		logger.Debugf("Served request for %s in %v", r.URL.Path, time.Since(begin))
		logger.Debugf("Body data field: %s", r.FormValue("data"))
	}
}

func MakeEchoMux(logger Logger) *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/", MakeHandler(EchoHandler, logger))
	return m
}

func ServeWithReload(m *http.ServeMux, logger Logger, c *comms) {
	for {
		addr := os.Getenv("ECHO_ADDR")
		if addr == "" {
			addr = "127.0.0.1:8080"
		}
		s := &http.Server{Addr: addr, Handler: m}

		go func(err chan<- error) {
			logger.Infof("Listening on %s", addr)
			if e := s.ListenAndServe(); e != nil && e != http.ErrServerClosed {
				err <- e
			}
			return
		}(c.err)

		select {
		case <-c.sig:
			logger.Info("SIGHUP recieved. Reloading...")
			begin := time.Now()
			logger.Debug("Halting Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := s.Shutdown(ctx); err != nil {
				logger.Debugf("Failed to gracefully shutdown server: %v", err)
			} else {
				logger.Debugf("Server halted in %v", time.Since(begin))
			}
		case <-c.done:
			logger.Info("Server shutting down...")
			begin := time.Now()
			logger.Debug("Halting Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := s.Shutdown(ctx); err != nil {
				logger.Debugf("Failed to gracefully shutdown server: %v", err)
			} else {
				logger.Debugf("Server halted in %v", time.Since(begin))
			}
			c.wg.Done()
			return
		}
	}
}

func RunServer(m *http.ServeMux, logger Logger) int {
	// First set up signal hanlding so that we can reload and stop.
	hups := make(chan os.Signal, 1)
	stop := make(chan os.Signal, 1)

	signal.Notify(hups, syscall.SIGHUP)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(hups)
	defer signal.Stop(stop)

	var waitgroup sync.WaitGroup

	done := make(chan struct{}, 1)
	err := make(chan error, 1)
	defer close(err)

	waitgroup.Add(1)

	go ServeWithReload(m, logger, &comms{&waitgroup, err, done, hups})

	exitCode := 1
	select {
	case <-stop:
		logger.Info("Interrupt/kill recieved. Exiting...")
		exitCode = 0
		close(done)
	case e := <-err:
		logger.Errorf("Encoutered error %v. Exiting...", e)
		close(done)
	}
	waitgroup.Wait()
	logger.Info("All waits done. Execution complete.")
	return exitCode
}

func RunEchoServer(logger Logger) int {
	m := MakeEchoMux(logger)
	return RunServer(m, logger)
}
