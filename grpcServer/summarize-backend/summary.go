package main

import (
	"context"

	pb "github.com/yseko789/grpcServer/proto"
)

func (grpcServer *grpcServer) ReadSummary(ctx context.Context, in *pb.Date) (*pb.Summary, error) {

	summaryData := &Summary{}

	res, err := collection.Doc(in.Date).Get(context.Background())
	if err != nil {
		return nil, err
	}

	err = res.DataTo(summaryData)
	if err != nil {
		return nil, err
	}

	return documentToSummary(in.Date, summaryData), nil
}
