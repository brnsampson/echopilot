package echoserver

import (
	"context"
	pb "github.com/brnsampson/echopilot/api/echo"
	"github.com/brnsampson/echopilot/pkg/echo"
	"google.golang.org/grpc"
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
func newGrpcServer(logger Logger, opts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterEchoServer(grpcServer, NewEchoServer(logger))
	return grpcServer
}
