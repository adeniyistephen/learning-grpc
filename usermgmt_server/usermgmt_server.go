package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/adeniyistephen/learning-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

// Create a new user
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Recerived: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(100))
	// Return user created
	return &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", port)
	// check for errors
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server register the UserManagementServer
	s := grpc.NewServer()
	// Register the UserManagementServer
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("Server listening on %v", lis.Addr())
	// Start the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
