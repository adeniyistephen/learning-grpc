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
	user_list *pb.UsersList
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		user_list: &pb.UsersList{},
	}
}


// Run the server
func (server *UserManagementServer) Run() error {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", port)
	// check for errors
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server register the UserManagementServer
	s := grpc.NewServer()
	// Register the UserManagementServer
	pb.RegisterUserManagementServer(s, server)
	log.Printf("Server listening on %v", lis.Addr())
	// Start the server
	return s.Serve(lis)
}

// Create a new user
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Recerived: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(100))
	// Return user created
	created_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   user_id,
	}
	// Add the user to the list
	s.user_list.Users = append(s.user_list.Users, created_user)
	return created_user, nil
}

// Get all users
func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UsersList, error) {
	// Return user created
	return server.user_list, nil
}

func main() {
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	// Run the server
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to run the server: %v", err)
	}
}
