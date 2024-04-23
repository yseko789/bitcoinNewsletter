package main

import (
	pb "github.com/yseko789/grpcServer/proto"
)

type Summary struct {
	Content string `firestore:"content"`
}

func documentToSummary(date string, data *Summary) *pb.Summary {
	return &pb.Summary{
		Date:    date,
		Content: data.Content,
	}
}
