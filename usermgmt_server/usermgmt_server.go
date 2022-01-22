package main

import (
	"context"
	"log"
	"fmt"
	"os"
	"net"

	"github.com/jackc/pgx/v4"
	pb "github.com/adeniyistephen/learning-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// UserManagemenetServer 
type UserManagementServer struct {
	conn *pgx.Conn
	first_user_creation bool
	pb.UnimplementedUserManagementServer
}

// NewUserManagementServer creates a new UserManagementServer Constructor
func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
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
	
	// create table if not exists
	createSql := `
	create table if not exists users(
		id SERIAL PRIMARY KEY,
		name text,
		age int
	);
	`

	// Execute the query
	_, err := s.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}
	
	s.first_user_creation = false

	log.Printf("Recerived: %v", in.GetName())

	// Create a new user
	created_user := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
	}

	// Begin a new transaction
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}

	// Insert the user
	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1,$2)",
		created_user.Name, created_user.Age)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}

	// Commit the transaction
	tx.Commit(context.Background())

	// Return the user created
	return created_user, nil
}

// Get all users
func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UsersList, error) {
	
	var users_list *pb.UsersList = &pb.UsersList{}
	// Quesry the database
	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Iterate over the rows
	for rows.Next() {
		// Create a new user
		user := pb.User{}
		// Scan the row into the user
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		// Append the user to the list
		users_list.Users = append(users_list.Users, &user)

	}
	// Return user created
	return users_list, nil
}

func main() {
	// Database connection string
	database_url := "postgres://postgres:mysecretpassword@localhost:5432/postgres"
	// sever instance
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	// Connect to the database
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())
	// Set the connection to the server
	user_mgmt_server.conn = conn
	user_mgmt_server.first_user_creation = true
	// Run the server
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to run the server: %v", err)
	}
}
