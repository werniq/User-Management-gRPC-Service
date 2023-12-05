package main

import (
	"fmt"
	pb "github.com/Chained/auth-service/github.com/Chained/auth-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var authServiceClient pb.AuthServiceClient
var authServiceAddress = fmt.Sprintf("localhost:%d", GrpcServerPort)

func initializeGRPCClient() error {
	/*
		creds, err := credentials.NewClientTLSFromFile("path/to/your/certificate.pem", "")
		if err != nil {
			return fmt.Errorf("failed to load TLS credentials: %v", err)
		}
	*/
	time.Sleep(time.Second * 5)

	conn, err := grpc.Dial(
		authServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	authServiceClient = pb.NewAuthServiceClient(conn)

	return nil
}
