package main

import (
	"context"
	protoinf "ez-cloud/hostpath-provisioner/grpc/inf"
	"fmt"
	"google.golang.org/grpc"
	"os"
)

func newNodeClient() {
	grpcClient, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stdout, "connect server error: %v", err)
		return
	}

	c := protoinf.NewNodeServiceClient(grpcClient)

	var r *protoinf.InfoResponse
	r, err = c.GetInfo(context.Background(), &protoinf.InfoRequest{})

	fmt.Fprintf(os.Stdout, "GetInfo ... [ %s => %s ]\n", r.GetId(), r.GetEndpoint())
}
