package main

import (
	"context"
	"fmt"
	"homework/internal/logger"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"homework/internal/server"
	pb "homework/pkg/files"
)

const (
	address         = ":50051"
	fileStoragePath = "./server_files"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Faild to create zap logger: %s", err)
	}
	logger.SetGlobal(zapLogger)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Errorf(ctx, "Failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()
	fh := server.NewHandler(fileStoragePath)
	pb.RegisterFileServiceServer(s, server.NewServer(fh))
	logger.Infof(ctx, "Server listening on %s", address)
	if err := s.Serve(lis); err != nil {
		logger.Errorf(ctx, "Failed to serve: %v\n", err)
		return
	}
}
