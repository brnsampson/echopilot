package echoserver

import (
	"context"
	"crypto/tls"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
type grpcRegistrant interface {
	RegisterWithServer(*grpc.Server) error
}

type serverOptions interface {
	getGrpcAddr() string
	getRestAddr() string
	getSrvOpts() []grpc.ServerOption
	getDialOpts() []grpc.DialOption
}

type server struct {
	logger     Logger
	registrant grpcRegistrant
	wg         sync.WaitGroup
	err        chan error
	done       chan struct{}
	stop       chan os.Signal
	sig        chan os.Signal
	exitCode   int
}

type serverOpts struct {
	grpcAddr string
	restAddr string
	srvOpts  []grpc.ServerOption
	dialOpts []grpc.DialOption
}

func NewServerOptsFromEnv(logger Logger) (*serverOpts, error) {
	var opts serverOpts

	// grpc Addr
	grpcAddr := os.Getenv("ECHO_GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = "127.0.0.1:8080"
		logger.Info("ECHO_GRPC_ADDR is empty. Defaulting to 127.0.0.1:8080")
	}
	opts.grpcAddr = grpcAddr

	// restAddr
	restAddr := os.Getenv("ECHO_REST_ADDR")
	if restAddr == "" {
		restAddr = "127.0.0.1:3000"
		logger.Info("ECHO_REST_ADDR is empty. Defaulting to 127.0.0.1:3000")
	}
	opts.restAddr = restAddr

	// srvOpts
	cert := os.Getenv("ECHO_SERVER_CERT")
	key := os.Getenv("ECHO_SERVER_KEY")

	if cert == "" {
		logger.Error("Error: ECHO_SERVER_CERT is empty when loading config!")
		return &opts, errors.New("NewServerOptionsFromEnv: ECHO_SERVER_CERT env var must be set")
	}

	if key == "" {
		logger.Error("Error: ECHO_SERVER_KEY is empty when loading config!")
		return &opts, errors.New("NewServerOptionsFromEnv: ECHO_SERVER_KEY env var must be set")
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		logger.Error("Error: Failed to create credentials from cert and key when loading config! Are ECHO_SERVER_CERT and ECHO_SERVER_KEY set properly?")
		return &opts, err
	}
	grpcCreds := grpc.Creds(creds)

	var srvOpts []grpc.ServerOption
	srvOpts = append(srvOpts, grpcCreds)

	opts.srvOpts = srvOpts

	// dialOpts
	var tlsSkipVerify bool
	skipVerify := os.Getenv("ECHO_GATEWAY_SKIP_VERIFY")
	if skipVerify == "" {
		logger.Info("ECHO_GATEWAY_SKIP_VERIFY is empty. Defaulting to true")
		tlsSkipVerify = true
	} else {
		tlsSkipVerify, err = strconv.ParseBool(skipVerify)
		if err != nil {
			logger.Info("ERROR: ECHO_GATEWAY_SKIP_VERIFY could not be parsed. Try true, false, or empty.")
			return &opts, err
		}
	}

	var dialOpts []grpc.DialOption
	if tlsSkipVerify {
		tlsSkipVerifyConfig := &tls.Config{
			InsecureSkipVerify: tlsSkipVerify,
		}
		tlsSkipVerifyOpt := grpc.WithTransportCredentials(credentials.NewTLS(tlsSkipVerifyConfig))
		dialOpts = append(dialOpts, tlsSkipVerifyOpt)
	} else {
		creds, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			logger.Error("Error: Failed to create client credentials from cert when loading config! Is ECHO_SERVER_CERT set properly?")
			return &opts, err
		}
		clientCreds := grpc.WithTransportCredentials(creds)
		dialOpts = append(dialOpts, clientCreds)
	}

	opts.dialOpts = dialOpts

	return &opts, nil
}

func (s *server) ServeWithReload() {
	for {
		var sopts serverOptions
		sopts, e := NewServerOptsFromEnv(s.logger)

		if e != nil {
			s.err <- e
		}

		grpcServer := grpc.NewServer(sopts.getSrvOpts()...)
		e = s.registrant.RegisterWithServer(grpcServer)

		if e != nil {
			s.err <- e
		}

		grpcListener, e := net.Listen("tcp", sopts.getGrpcAddr())
		if e != nil {
			s.err <- e
		}

		ctx, cancel := context.WithCancel(context.Background())
		gm, e := NewGatewayMux(ctx, sopts.getGrpcAddr(), sopts.getDialOpts())
		if e != nil {
			s.err <- e
		}

		httpServ := &http.Server{Addr: sopts.getRestAddr(), Handler: gm}

		// Start GRPC server
		go func(err chan<- error) {
			s.logger.Infof("GRPC server listening on %s", sopts.getGrpcAddr())
			if e := grpcServer.Serve(grpcListener); e != nil {
				err <- e
			}
			return
		}(s.err)

		// Start http server for swagger and grpc rest gateway
		go func(err chan<- error) {
			s.logger.Infof("Rest gateway listening on %s", sopts.getRestAddr())
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

func NewServer(logger Logger) *server {
	// First set up signal hanlding so that we can reload and stop.
	hups := make(chan os.Signal, 1)
	stop := make(chan os.Signal, 1)

	signal.Notify(hups, syscall.SIGHUP)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var waitgroup sync.WaitGroup

	done := make(chan struct{}, 1)
	err := make(chan error, 1)

	registrant, e := NewServerRegistrant(logger)

	if e != nil {
		logger.Error("Error: failed to create GRPC server registrant with call to NewServerRegistrant()")
		err <- e
	}

	// If exitCode is -1 then execution has not completed yet.
	server := &server{logger, registrant, waitgroup, err, done, stop, hups, -1}
	return server
}
