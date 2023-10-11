package echo

import (
    "net/http"
    "github.com/brnsampson/echopilot/rpc/echo/internal"
	pb "github.com/brnsampson/echopilot/proto/gen/echo/v1"
    "github.com/charmbracelet/log"
)

// Requests
func NewStringRequest(req string) *pb.EchoStringRequest {
    return &pb.EchoStringRequest{Content: req}
}

func NewIntRequest(req int32) *pb.EchoIntRequest {
    return &pb.EchoIntRequest{Content: req}
}

func ReadStringResult(res *pb.EchoStringResponse) string {
    return res.Content
}

func ReadIntResult(res *pb.EchoIntResponse) int32 {
    return res.Content
}

// Service
func NewService(log *log.Logger) Service {

	return Service{ log.With("service", "echo") }
}

type Service struct{
    *log.Logger
}

func (s *Service) EchoString(req *pb.EchoStringRequest) (*pb.EchoStringResponse, error) {
	res := &pb.EchoStringResponse{Content: req.Content}
	return res, nil
}

func (s *Service) EchoInt(req *pb.EchoIntRequest) (*pb.EchoIntResponse, error) {
	res := &pb.EchoIntResponse{Content: req.Content}
	return res, nil
}

func (s *Service) GetHandler() (string, http.Handler) {
    connect := internal.NewEchoConnectServer(s)
    return connect.GetHttpHandler()
}
