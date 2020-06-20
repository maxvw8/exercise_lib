package main

import (
	"context"
	"fmt"
	"log"

	pbexrs "github.com/maxvw8/exercise_lib/pbexrs/v1"
	grpc "google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting Blog Service Client")

	// Establish insecure grpc options (no TLS)
	requestOpts := grpc.WithInsecure()
	// Dial the server, returns a client connection
	conn, err := grpc.Dial("localhost:10000", requestOpts)
	if err != nil {
		log.Fatalf("Unable to establish client connection to localhost:10000: %v", err)
	}
	defer conn.Close()
	// Instantiate the BlogServiceClient with our client connection to the server
	client := pbexrs.NewExerciseServiceClient(conn)
	ctx := context.Background()
	list, err := client.ListExercises(ctx, &pbexrs.ListExercisesRequest{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("list:%v", list)
}
