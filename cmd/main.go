package main

import (
	pb "github.com/Chained/auth-service/github.com/Chained/auth-service"
	"github.com/Chained/auth-service/internal/models"
	"os"
	"sync"
)

type Application struct {
	DB     *models.DatabaseModel
	Client pb.AuthServiceClient
	Server *pb.AuthServiceServer
}

var app *Application

func init() {
	var err error

	// Initialize gRPC client connection during startup
	go func() {
		err = initializeGRPCClient()
		if err != nil {
			ErrorLogger.Printf("initializing grpc client")
		}
		app.Client = authServiceClient
	}()

	if err != nil {
		ErrorLogger.Printf("Failed to initialize gRPC client: %v", err)
		os.Exit(1)
	}

	app.DB, err = models.NewDBModel()
	if err != nil {
		ErrorLogger.Printf("Establishing database connection: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	var wg sync.WaitGroup

	// HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		startHTTPServer()
	}()

	// gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		startGRPCServer()
	}()

	// wait for both servers to finish
	wg.Wait()
}
