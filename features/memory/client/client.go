package echo

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	pb "github.com/brnsampson/echopilot/proto/gen/echo/v1"
	"github.com/brnsampson/echopilot/proto/gen/echo/v1/echov1connect"
	"github.com/bufbuild/connect-go"
    "github.com/brnsampson/echopilot/pkg/option"
)

type EchoClient interface {
    // This should match the service methods exposed in the protobuf service
    EchoString(request pb.EchoStringRequest) (pb.EchoStringResponse, error)
    EchoInt(request pb.EchoIntRequest) (pb.EchoIntResponse, error)
}

type RemoteEchoClient struct {
	connectClient echov1connect.EchoServiceClient
}

func (ec *RemoteEchoClient) EchoString(request *pb.EchoStringRequest) (*pb.EchoStringResponse, error) {
	req := connect.NewRequest(request)
	response, err := ec.connectClient.EchoString(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return response.Msg, nil
}

func (ec *RemoteEchoClient) EchoInt(request *pb.EchoIntRequest) (*pb.EchoIntResponse, error) {
	req := connect.NewRequest(request)
	response, err := ec.connectClient.EchoInt(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return response.Msg, nil
}

func NewRemoteEchoClient(addr string, timeout option.Option[time.Duration], skipVerify option.Option[bool]) (*RemoteEchoClient, error) {
	if addr == "" {
		addr = "127.0.0.1:3000"
	}

	sv := skipVerify.UnwrapOrDefault(false)

	tlsConf := tls.Config{InsecureSkipVerify: sv}
	transport := http.Transport{TLSClientConfig: &tlsConf}

	to := timeout.UnwrapOrDefault(time.Duration(10) * time.Second)

	client := http.Client{Timeout: to, Transport: &transport}

	echoclient := echov1connect.NewEchoServiceClient(&client, addr)

	return &RemoteEchoClient{echoclient}, nil
}
