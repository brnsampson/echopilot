package echoclient

import (
	"context"
	"fmt"
	pb "github.com/brnsampson/echopilot/proto/echo"
	"google.golang.org/grpc"
	"os"
	"time"
)

func EchoString(client pb.EchoClient, req *pb.EchoRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	response, err := client.Echo(ctx, req)
	if err != nil {
		fmt.Printf("%v.Echo(_) = _, %v: ", client, err)
	}
	fmt.Println(response)
}

func GetEcho() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	addr := os.Getenv("ECHO_GRPC_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		fmt.Println("Error while dialing")
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	client := pb.NewEchoClient(conn)
	EchoString(client, &pb.EchoRequest{Content: "tester"})
}
