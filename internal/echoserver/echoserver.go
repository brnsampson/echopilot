package echoserver

import (
	"context"
	"crypto/tls"
	"errors"
	pb "github.com/brnsampson/echopilot/api/echo"
	"github.com/brnsampson/echopilot/pkg/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"strconv"
	"time"
)

// Application specific functionas and behavior.

// Begin grpc server for echo logic
type echoServer struct {
	pb.UnimplementedEchoServer
	logger Logger
}

func NewEchoServer(logger Logger) *echoServer {
	return &echoServer{logger: logger}
}

func (es *echoServer) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoReply, error) {
	es.logger.Debug("Recieved Echo request: ", ctx)

	begin := time.Now()
	result, err := echo.EchoString(req.Content)
	if err != nil {
		return &pb.EchoReply{Content: ""}, err
	}
	es.logger.Infof("Echo request handled in %v", time.Since(begin))
	es.logger.Debugf("Replying to Echo request with: %s", result)
	return &pb.EchoReply{Content: result}, nil
}

// End grpc server for echo logic

// Generic wrapper for grpc server to be called in server.go
func NewServerRegistrant(logger Logger) (grpcRegistrant, error) {
	return NewEchoServer(logger), nil
}

func (es *echoServer) RegisterWithServer(grpcServer *grpc.Server) error {
	pb.RegisterEchoServer(grpcServer, es)
	return nil
}

//Implement ServerOptions interface
type serverOpts struct {
	grpcAddr string
	restAddr string
	srvOpts  []grpc.ServerOption
	dialOpts []grpc.DialOption
}

func (so *serverOpts) getGrpcAddr() string {
	return so.grpcAddr
}

func (so *serverOpts) getRestAddr() string {
	return so.restAddr
}

func (so *serverOpts) getSrvOpts() []grpc.ServerOption {
	return so.srvOpts
}

func (so *serverOpts) getDialOpts() []grpc.DialOption {
	return so.dialOpts
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
