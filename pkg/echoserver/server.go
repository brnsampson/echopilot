package echoserver

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

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
	logger     Logger
	grpcServer *grpc.Server
	s_opts     []grpc.ServerOption
	d_opts     []grpc.DialOption
	wg         sync.WaitGroup
	err        chan error
	done       chan struct{}
	stop       chan os.Signal
	sig        chan os.Signal
	exitCode   int
}

func (s *server) ServeWithReload() {
	for {
		grpc_addr := os.Getenv("ECHO_GRPC_ADDR")
		if grpc_addr == "" {
			grpc_addr = "127.0.0.1:8080"
		}

		http_addr := os.Getenv("ECHO_HTTP_ADDR")
		if http_addr == "" {
			http_addr = "127.0.0.1:3000"
		}

		grpc_listener, e := net.Listen("tcp", grpc_addr)
		if e != nil {
			s.err <- e
		}

		ctx, cancel := context.WithCancel(context.Background())
		gm, e := NewGatewayMux(ctx, grpc_addr, s.d_opts)
		if e != nil {
			s.err <- e
		}
		http_serv := &http.Server{Addr: http_addr, Handler: gm}

		// Start http server for swagger and grpc rest gateway
		go func(err chan<- error) {
			s.logger.Infof("Rest gateway listening on %s", http_addr)
			defer cancel()
			if e := http_serv.ListenAndServe(); e != nil && e != http.ErrServerClosed {
				err <- e
			}
		}(s.err)

		// Start GRPC server
		go func(err chan<- error) {
			s.logger.Infof("GRPC server listening on %s", grpc_addr)
			if e := s.grpcServer.Serve(grpc_listener); e != nil {
				err <- e
			}
			return
		}(s.err)

		select {
		case <-s.sig:
			s.logger.Info("SIGHUP recieved. Reloading...")
			begin := time.Now()
			s.logger.Debug("Halting HTTP Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := http_serv.Shutdown(ctx); err != nil {
				s.logger.Debugf("Failed to gracefully shutdown server: %v", err)
			} else {
				s.logger.Debugf("HTTP server halted in %v", time.Since(begin))
			}

			s.logger.Debug("Halting GRPC Server...")
			begin = time.Now()
			s.grpcServer.GracefulStop()
			s.logger.Debugf("GRPC server halted in %v", time.Since(begin))
		case <-s.done:
			s.logger.Info("Server shutting down...")
			begin := time.Now()
			s.logger.Debug("Halting HTTP Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := http_serv.Shutdown(ctx); err != nil {
				s.logger.Debugf("Failed to gracefully shutdown HTTP server: %v", err)
			} else {
				s.logger.Debugf("HTTP server halted in %v", time.Since(begin))
			}

			s.logger.Debug("Halting GRPC Server...")
			begin = time.Now()
			s.grpcServer.GracefulStop()
			s.logger.Debugf("GRPC server halted in %v", time.Since(begin))
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

	var srv_opts []grpc.ServerOption
	gs := newGrpcServer(logger, srv_opts)

	var dial_opts []grpc.DialOption
	dial_opts = append(dial_opts, grpc.WithInsecure())

	// If exitCode is -1 then execution has not completed yet.
	server := &server{logger, gs, srv_opts, dial_opts, waitgroup, err, done, stop, hups, -1}
	return server
}
