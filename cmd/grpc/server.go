package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/maxvw8/exercise_lib/exrs"
	"github.com/maxvw8/exercise_lib/exrs/storage/mongodb"
	pbexrs "github.com/maxvw8/exercise_lib/pbexrs/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	fmt.Println("Starting server on port :50051...")
	// Start our listener, 50051 is the default gRPC port
	listener, err := net.Listen("tcp", ":50051")
	// Handle errors if any
	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}
	// Set options, here we can configure things like TLS support
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	}
	// Create new gRPC server with (blank) options
	s := grpc.NewServer(opts...)
	//create repository connection
	repo, err := mongodb.New("myDB")
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	// Create BlogService type
	srv, err := exrs.Server(repo)
	if err != nil {
		log.Fatal(err)
	}
	// Register the service with the server
	pbexrs.RegisterExerciseServiceServer(s, srv)
	// Start the server in a child routine
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50051")

	// Right way to stop the server using a SHUTDOWN HOOK

	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed
	// And our main routine keeps running
	// If the main routine were to shutdown so would the child routine that is Serving the server
	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	fmt.Println("Done.")
}
