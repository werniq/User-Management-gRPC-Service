package main

import (
	"fmt"
	pb "github.com/Chained/auth-service/github.com/Chained/auth-service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
	"os"
)

// startHTTPServer configures routes for gin.Engine and runs server on the port 8080
func startHTTPServer() {
	r := gin.Default()
	app.SetupRoutes(r)

	if err := r.Run(fmt.Sprintf(":%d", HTTPServerPort)); err != nil {
		ErrorLogger.Printf("running HTTP server on port %d: %v\n", 8080, err)
	}
}

// startGRPCServer initializes and starts the gRPC server.
// It listens on the specified port for incoming gRPC connections,
// registers the AuthServiceServer implementation, and serves requests.
// The function logs server initialization information and any errors
// encountered during the server's lifecycle.
func startGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", GrpcServerPort))
	if err != nil {
		ErrorLogger.Printf("failed to listen: %v\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &UserServiceServer{})

	InfoLogger.Printf("gRPC server is listening on port %d...\n", GrpcServerPort)

	if err := s.Serve(lis); err != nil {
		ErrorLogger.Printf("failed to serve: %v\n", err)
	}
}
