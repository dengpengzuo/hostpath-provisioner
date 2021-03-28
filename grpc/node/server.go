package main

import (
	"context"
	protoinf "ez-cloud/hostpath-provisioner/grpc/inf"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"time"
)

type node struct {
}

func (s *node) GetInfo(context.Context, *protoinf.InfoRequest) (*protoinf.InfoResponse, error) {
	return &protoinf.InfoResponse{Id: "1", Endpoint: "tcp://localhost:2181"}, nil
}

func (s *node) Health(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func newNodeServer() {
	grpcServer := grpc.NewServer(grpc.ConnectionTimeout(30 * time.Second))

	nodeService := &node{}
	protoinf.RegisterNodeServiceServer(grpcServer, nodeService)

	lis, _ := net.Listen("tcp", "localhost:8080")
	grpcServer.Serve(lis)
}
