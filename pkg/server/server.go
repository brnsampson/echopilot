package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

    "github.com/charmbracelet/log"
)

type ServerOptions interface {
	GetAddr(bool) (string, error)
	GetTlsConfig(bool) (*tls.Config, error)
	GetTlsEnabled(bool) (bool, error)
}

type Server struct {
	logger   *log.Logger
	wg       *sync.WaitGroup
	err      chan error
	done     chan struct{}
	stop     chan os.Signal
	reload   chan os.Signal
	exitCode int
}

func (s *Server) ServeWithReload(router http.Handler, sopts ServerOptions) {
	for {
		addr, err := sopts.GetAddr(true)
		if err != nil {
			s.logger.Error("Error: failed to refresh server config. May use some old settings.")
		}

		tlsConf, _ := sopts.GetTlsConfig(false)
		if err != nil {
			s.logger.Error("Error: failed to refresh server config. May use some old settings.")
		}

		tlsEnabled, _ := sopts.GetTlsEnabled(false)
		if err != nil {
			s.logger.Error("Error: failed to refresh server config. May use some old settings.")
		}

        stdlog := s.logger.StandardLog(log.StandardLogOptions{
            ForceLevel: log.ErrorLevel,
        })
		httpServ := &http.Server{
			Addr:         addr,
			Handler:      router,
            ErrorLog:     stdlog,
			TLSConfig:    tlsConf,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		go func(err chan<- error) {
			if tlsEnabled {
				// Note that the certificate is already embedded in the tlsConf and that will override
				// any cert/key filenames we pass anyways.
				s.logger.Infof("https server listening on %s", addr)
				if e := httpServ.ListenAndServeTLS("", ""); e != nil && e != http.ErrServerClosed {
					err <- e
				}
			} else {
				s.logger.Infof("http server listening on %s", addr)
				if e := httpServ.ListenAndServe(); e != nil && e != http.ErrServerClosed {
					err <- e
				}
			}
		}(s.err)

		select {
		case <-s.reload:
			s.logger.Info("SIGHUP received. Reloading...")
			begin := time.Now()
			s.logger.Debug("Halting HTTP Server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := httpServ.Shutdown(ctx); err != nil {
				s.logger.Debugf("Failed to gracefully shutdown server: %v", err)
			} else {
				s.logger.Debugf("HTTP server halted in %v", time.Since(begin))
			}
			cancel()
		case <-s.done:
			s.logger.Info("Server shutting down...")
			begin := time.Now()
			s.logger.Debug("Halting HTTP Server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := httpServ.Shutdown(ctx); err != nil {
				s.logger.Debugf("Failed to gracefully shutdown HTTP server: %v", err)
			} else {
				s.logger.Debugf("HTTP server halted in %v", time.Since(begin))
			}

			s.wg.Done()
			cancel()
			return
		}
	}
}

func (s *Server) Run(router http.Handler, sopts ServerOptions) {
	s.wg.Add(1)
	go s.waitOnInterrupt()
	go s.ServeWithReload(router, sopts)
	return
}

func (s *Server) BlockingRun(router http.Handler, opts ServerOptions) int {
	s.wg.Add(1)
	go s.ServeWithReload(router, opts)
	s.waitOnInterrupt()
	return s.exitCode
}

func (s *Server) waitOnInterrupt() {
	select {
	case <-s.stop:
		s.logger.Info("Interrupt/kill received. Exiting...")
		s.exitCode = 0
		s.Shutdown()
	case e := <-s.err:
		s.logger.Errorf("Encountered error %v. Exiting...", e)
		s.Shutdown()
		s.exitCode = 1
	}
}

func (s *Server) Shutdown() {
	signal.Stop(s.reload)
	signal.Stop(s.stop)
	close(s.done)

	s.wg.Wait()
	s.logger.Info("All waits done. Server execution complete.")

	close(s.stop)
	close(s.reload)
	close(s.err)
	return
}

func NewServer(logger *log.Logger) *Server {
	// First set up signal handling so that we can reload and stop.
	hups := make(chan os.Signal, 1)
	stop := make(chan os.Signal, 1)

	signal.Notify(hups, syscall.SIGHUP)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var waitgroup sync.WaitGroup

	done := make(chan struct{}, 1)
	err := make(chan error, 1)

	// If exitCode is -1 then execution has not completed yet.
	server := &Server{logger.With("package", "server"), &waitgroup, err, done, stop, hups, -1}
	return server
}
