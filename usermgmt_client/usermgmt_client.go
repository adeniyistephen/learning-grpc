package main

import (
	"context"
	"log"
	"time"

	pb "github.com/adeniyistephen/learning-grpc/usermgmt"
	"google.golang.org/grpc"
)

// address for the client to connect to the server
const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Allow the client to access the server connection
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new user
	var new_users = make(map[string]int32)
	new_users["Alice"] = 43
	new_users["Bob"] = 30
	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		log.Printf(`User Details:
		NAME: %s
		AGE: %d
		ID: %d`, r.GetName(), r.GetAge(), r.GetId())

	}

	// Get Users Client
	params := &pb.GetUsersParams{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	log.Print("\nUSER LIST: \n")
	// Print the user list
	for _, user := range r.GetUsers() {
		log.Printf(`User Details:
		NAME: %s
		AGE: %d
		ID: %d`, user.GetName(), user.GetAge(), user.GetId())
	}
}
