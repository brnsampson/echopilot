package internal

import (
	"context"
	"net/http"

	pb "github.com/brnsampson/echopilot/proto/gen/echo/v1"
	"github.com/brnsampson/echopilot/proto/gen/echo/v1/echov1connect"
    "connectrpc.com/connect"
)

type EchoService interface {
    EchoString(req *pb.EchoStringRequest) (*pb.EchoStringResponse, error)
    EchoInt(req *pb.EchoIntRequest) (*pb.EchoIntResponse, error)
}

// EchoConnectServer implements methods to handle requests by calling service methods.
// The logic here should be minimal and most work should be handled either in middleware
// or the actual service pkg. Extensive logic here indicates that our API and domain
// logic are drifting.
func NewEchoConnectServer(service EchoService) *EchoConnectServer {
	return &EchoConnectServer{service}
}

type EchoConnectServer struct {
	EchoService
}

func (es *EchoConnectServer) GetHttpHandler() (string, http.Handler) {
	return echov1connect.NewEchoServiceHandler(es)
}

func (es *EchoConnectServer) EchoString(ctx context.Context, req *connect.Request[pb.EchoStringRequest]) (*connect.Response[pb.EchoStringResponse], error) {

	result, err := es.EchoService.EchoString(req.Msg)
	if err != nil {
		res := connect.NewResponse(&pb.EchoStringResponse{Content: ""})
		return res, err
	}

	res := connect.NewResponse(result)
	return res, nil
}

func (es *EchoConnectServer) EchoInt(ctx context.Context, req *connect.Request[pb.EchoIntRequest]) (*connect.Response[pb.EchoIntResponse], error) {

	result, err := es.EchoService.EchoInt(req.Msg)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(result)
	return res, nil
}
