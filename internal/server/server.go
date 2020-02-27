package server

import (
	"context"
	domainserver "github.com/brnsampson/echopilot/internal/echoserver"
	"github.com/brnsampson/echopilot/pkg/logger"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Common server implementation boilerplate
type grpcRegistrant interface {
	RegisterWithServer(*grpc.Server) error
}

type serverOptions interface {
	GetGrpcAddr() string
	GetRestAddr() string
	GetSrvOpts() []grpc.ServerOption
	GetDialOpts() []grpc.DialOption
}

type server struct {
	logger     logger.Logger
	flags      *pflag.FlagSet
	registrant grpcRegistrant
	wg         sync.WaitGroup
	err        chan error
	done       chan struct{}
	stop       chan os.Signal
	sig        chan os.Signal
	exitCode   int
}

func (s *server) ServeWithReload() {
	for {
		var sopts serverOptions
		sopts, e := domainserver.NewServerOpts(s.logger, s.flags)

		if e != nil {
			s.err <- e
		}

		grpcServer := grpc.NewServer(sopts.GetSrvOpts()...)
		e = s.registrant.RegisterWithServer(grpcServer)

		if e != nil {
			s.err <- e
		}

		grpcListener, e := net.Listen("tcp", sopts.GetGrpcAddr())
		if e != nil {
			s.err <- e
		}

		ctx, cancel := context.WithCancel(context.Background())
		gm, e := NewGatewayMux(ctx, sopts.GetGrpcAddr(), sopts.GetDialOpts())
		if e != nil {
			s.err <- e
		}

		httpServ := &http.Server{Addr: sopts.GetRestAddr(), Handler: gm}

		// Start GRPC server
		go func(err chan<- error) {
			s.logger.Infof("GRPC server listening on %s", sopts.GetGrpcAddr())
			if e := grpcServer.Serve(grpcListener); e != nil {
				err <- e
			}
			return
		}(s.err)

		// Start http server for swagger and grpc rest gateway
		go func(err chan<- error) {
			s.logger.Infof("Rest gateway listening on %s", sopts.GetRestAddr())
			defer cancel()
			if e := httpServ.ListenAndServe(); e != nil && e != http.ErrServerClosed {
				err <- e
			}
		}(s.err)

		select {
		case <-s.sig:
			s.logger.Info("SIGHUP recieved. Reloading...")
			begin := time.Now()
			s.logger.Debug("Halting HTTP Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := httpServ.Shutdown(ctx); err != nil {
				s.logger.Debugf("Failed to gracefully shutdown server: %v", err)
			} else {
				s.logger.Debugf("HTTP server halted in %v", time.Since(begin))
			}

			s.logger.Debug("Halting GRPC Server...")
			begin = time.Now()
			grpcServer.GracefulStop()
			s.logger.Debugf("GRPC server halted in %v", time.Since(begin))
		case <-s.done:
			s.logger.Info("Server shutting down...")
			begin := time.Now()
			s.logger.Debug("Halting HTTP Server...")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if err := httpServ.Shutdown(ctx); err != nil {
				s.logger.Debugf("Failed to gracefully shutdown HTTP server: %v", err)
			} else {
				s.logger.Debugf("HTTP server halted in %v", time.Since(begin))
			}

			s.logger.Debug("Halting GRPC Server...")
			begin = time.Now()
			grpcServer.GracefulStop()
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

func NewServer(logger logger.Logger, flags *pflag.FlagSet) *server {
	// First set up signal hanlding so that we can reload and stop.
	hups := make(chan os.Signal, 1)
	stop := make(chan os.Signal, 1)

	signal.Notify(hups, syscall.SIGHUP)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var waitgroup sync.WaitGroup

	done := make(chan struct{}, 1)
	err := make(chan error, 1)

	registrant, e := domainserver.NewServerRegistrant(logger)

	if e != nil {
		logger.Error("Error: failed to create GRPC server registrant with call to NewServerRegistrant()")
		err <- e
	}

	// If exitCode is -1 then execution has not completed yet.
	server := &server{logger, flags, registrant, waitgroup, err, done, stop, hups, -1}
	return server
}
