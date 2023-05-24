package main

import (
	"fmt"
	"github.com/zondax/keyringPoc/keyring"
	"github.com/zondax/keyringPoc/keyring/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	fmt.Println("HELLo")
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := server.NewApiServer()
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	keyring.RegisterKeyringServiceServer(grpcServer, s)
	grpcServer.Serve(lis)
}
