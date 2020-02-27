package echoclient

import (
	"context"
	"crypto/tls"
	pb "github.com/brnsampson/echopilot/api/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"strconv"
	"time"
)

func EchoString(client pb.EchoClient, req *pb.EchoRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.Echo(ctx, req)
	if err != nil {
		return "", err
	}
	return response.Content, nil
}

func GetEcho(request string) (string, error) {
	var opts []grpc.DialOption
	addr := os.Getenv("ECHO_GRPC_ADDR")

	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	var skip bool
	var err error
	skipVerify := os.Getenv("ECHO_CLIENT_SKIP_VERIFY")
	if skipVerify == "" {
		skip = false
	} else {
		skip, err = strconv.ParseBool(skipVerify)
		if err != nil {
			return "", err
		}
	}

	tlsConf := &tls.Config{InsecureSkipVerify: skip}
	tlsOpt := credentials.NewTLS(tlsConf)
	opts = append(opts, grpc.WithTransportCredentials(tlsOpt))
	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		return "", err
	}

	defer conn.Close()

	client := pb.NewEchoClient(conn)
	result, err := EchoString(client, &pb.EchoRequest{Content: request})

	if err != nil {
		return "", err
	}

	return result, nil
}
