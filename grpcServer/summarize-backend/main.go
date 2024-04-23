package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	pb "github.com/yseko789/grpcServer/proto"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.SummaryServiceServer
}

var collection *firestore.CollectionRef

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("projectID"))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	logger.Info("firestore connection pool established")
	collection = client.Collection("summaries")

	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error listening")
	}
	grpcS := grpc.NewServer()
	pb.RegisterSummaryServiceServer(grpcS, &grpcServer{})
	if err := grpcS.Serve(list); err != nil {
		log.Fatal(err)
	}
}
