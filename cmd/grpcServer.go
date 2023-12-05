package main

import (
	"context"
	pb "github.com/Chained/auth-service/github.com/Chained/auth-service"
	"log"
	"os"
)

type UserServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

var (
	ErrorLogger    = log.New(os.Stdout, "\x1b[31m[ERROR]\x1b[0m : \t", log.Lshortfile|log.Ldate|log.Ltime)
	InfoLogger     = log.New(os.Stdout, "[INFO]: \t", log.Ldate|log.Ltime)
	GrpcServerPort = 50051
	HTTPServerPort = 8080
)

// Authorize function hashes user.Password and inserts it into db
func (s *UserServiceServer) Authorize(ctx context.Context, user *pb.User) (*pb.User, error) {
	var err error

	err = app.DB.CreateUser(
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	)

	if err != nil {
		return nil, err
	}

	InfoLogger.Println("User successfully created.")

	return user, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, gr *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := app.DB.GetUser(int(gr.Id))
	if err != nil {
		return nil, err
	}

	InfoLogger.Println("User successfully retrieved: ", user)

	return &pb.GetUserResponse{User: user}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, ur *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	err := app.DB.UpdateUser(ur)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		Status: "Successfully updated",
	}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, ur *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	var ures *pb.DeleteUserResponse

	err := app.DB.DeleteUser(int(ur.Id))
	if err != nil {
		// ures.ResponseCode = 0
		return ures, err
	}

	ures.ResponseCode = 1
	return ures, nil
}
