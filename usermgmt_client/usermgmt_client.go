package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "github.com/adeniyistephen/learning-grpc/usermgmt"
	"google.golang.org/grpc"
)

// address for the client to connect to the server
const (
	address = "localhost:50051"
)

// Creating out client from the generated proto file
var c pb.UserManagementClient

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Allow the client to access the server connection
	c = pb.NewUserManagementClient(conn)

	http.HandleFunc("/create", createUser)
	http.HandleFunc("/get", getUser)
	log.Fatal(http.ListenAndServe(":4040", nil))
}

// Create a new user handler with http/net go library
func createUser(w http.ResponseWriter, r *http.Request) {

	// Handling context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create a new user
	var new_users = make(map[string]int32)
	new_users["Alice"] = 43
	new_users["Bob"] = 30
	for name, age := range new_users {
		res, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		log.Printf(`User Details:
		NAME: %s
		AGE: %d
		ID: %d`, res.GetName(), res.GetAge(), res.GetId())

		// Print the user created
		fmt.Fprint(w, res)
	}
}

// Get users handler with http/net go library
func getUser(w http.ResponseWriter, r *http.Request) {
	// Handling context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	// Get Users Client
	params := &pb.GetUsersParams{}
	res, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}

	log.Print("\nUSER LIST: \n")

	// Print the user list
	for _, user := range res.GetUsers() {
		log.Printf(`User Details:
		NAME: %s
		AGE: %d
		ID: %d`, user.GetName(), user.GetAge(), user.GetId())
	}

	// Print the user list
	fmt.Fprint(w, res)
}
