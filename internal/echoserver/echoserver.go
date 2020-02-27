package echoserver

import (
	"context"
	"crypto/tls"
	"errors"
	pb "github.com/brnsampson/echopilot/api/echo"
	"github.com/brnsampson/echopilot/pkg/echo"
	"github.com/brnsampson/echopilot/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

// Application specific functionas and behavior.

// rumtime mux for REST gateway
func NewRuntimeMux(ctx context.Context, grpcAddr string, opts []grpc.DialOption) (*runtime.ServeMux, error) {
	rmux := runtime.NewServeMux()
	err := pb.RegisterEchoHandlerFromEndpoint(ctx, rmux, grpcAddr, opts)
	if err != nil {
		return rmux, err
	}
	return rmux, nil
}

// Begin grpc server for echo logic
type echoServer struct {
	pb.UnimplementedEchoServer
	logger logger.Logger
}

func NewEchoServer(logger logger.Logger) *echoServer {
	return &echoServer{logger: logger}
}

func (es *echoServer) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoReply, error) {
	es.logger.Debug("Recieved Echo request: ", ctx)

	begin := time.Now()
	result, err := echo.EchoString(req.Content)
	if err != nil {
		return &pb.EchoReply{Content: ""}, err
	}
	es.logger.Infof("Echo request handled in %+v", time.Since(begin))
	es.logger.Debugf("Replying to Echo request with: %s", result)
	return &pb.EchoReply{Content: result}, nil
}

// End grpc server for echo logic

// Generic wrapper for grpc server to be called in server.go
func NewServerRegistrant(logger logger.Logger) (*echoServer, error) {
	return NewEchoServer(logger), nil
}

func (es *echoServer) RegisterWithServer(grpcServer *grpc.Server) error {
	pb.RegisterEchoServer(grpcServer, es)
	return nil
}

//Implement ServerOptions interface
type ServerOpts struct {
	grpcAddr string
	restAddr string
	srvOpts  []grpc.ServerOption
	dialOpts []grpc.DialOption
}

func (so *ServerOpts) GetGrpcAddr() string {
	return so.grpcAddr
}

func (so *ServerOpts) GetRestAddr() string {
	return so.restAddr
}

func (so *ServerOpts) GetSrvOpts() []grpc.ServerOption {
	return so.srvOpts
}

func (so *ServerOpts) GetDialOpts() []grpc.DialOption {
	return so.dialOpts
}

func NewServerOpts(logger logger.Logger, flags *pflag.FlagSet) (*ServerOpts, error) {
	var opts ServerOpts

	conf, err := NewFullConfig(logger, flags)
	if err != nil {
		logger.Error("Error: could not load config!")
		return &opts, err
	}

	logger.Infof("Creating server opts from merged config: %+v", conf)

	opts.grpcAddr = conf.GrpcAddress

	opts.restAddr = conf.RestAddress

	if conf.TlsCert == "" {
		logger.Error("Error: TlsCert is empty when loading config!")
		return &opts, errors.New("NewServerOptionsFromEnv: tls certificate location must be set")
	}

	if conf.TlsKey == "" {
		logger.Error("Error: tlsKey is empty when loading config!")
		return &opts, errors.New("NewServerOptionsFromEnv: tls key location must be set")
	}

	creds, err := credentials.NewServerTLSFromFile(conf.TlsCert, conf.TlsKey)
	if err != nil {
		logger.Error("Error: Failed to create credentials from cert and key when loading config! Are the tls cert and key paths set properly?")
		return &opts, err
	}
	grpcCreds := grpc.Creds(creds)

	var srvOpts []grpc.ServerOption
	srvOpts = append(srvOpts, grpcCreds)

	opts.srvOpts = srvOpts

	var dialOpts []grpc.DialOption
	if conf.TlsSkipVerify {
		tlsSkipVerifyConfig := &tls.Config{
			InsecureSkipVerify: conf.TlsSkipVerify,
		}
		tlsSkipVerifyOpt := grpc.WithTransportCredentials(credentials.NewTLS(tlsSkipVerifyConfig))
		dialOpts = append(dialOpts, tlsSkipVerifyOpt)
	} else {
		creds, err := credentials.NewClientTLSFromFile(conf.TlsCert, "")
		if err != nil {
			logger.Error("Error: Failed to create client credentials from cert when loading config! Is the tls certificate path set properly?")
			return &opts, err
		}
		clientCreds := grpc.WithTransportCredentials(creds)
		dialOpts = append(dialOpts, clientCreds)
	}

	opts.dialOpts = dialOpts

	return &opts, nil
}
