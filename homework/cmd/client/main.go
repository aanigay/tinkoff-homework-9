package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"homework/internal/logger"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"homework/internal/client"
	pb "homework/pkg/files"
)

const (
	address = "localhost:50051"
	timeout = 1 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Faild to create zap logger: %s", err)
	}
	logger.SetGlobal(zapLogger)
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)

	uploadFilename := uploadCmd.String("file", "", "File to upload")
	infoFilename := infoCmd.String("file", "", "File to get information")

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  client upload --file <filename>")
		fmt.Println("  client list")
		fmt.Println("  client info --file <filename>")
		os.Exit(1)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	c := pb.NewFileServiceClient(conn)
	switch os.Args[1] {
	case "upload":
		err = uploadCmd.Parse(os.Args[2:])
		if err != nil {
			logger.Errorf(ctx, "Failed to parse command, please try again")
			os.Exit(1)
		}
		if *uploadFilename == "" {
			logger.Errorf(ctx, "Please provide a filename to upload.")
			os.Exit(1)
		}
		err = client.UploadFile(ctx, c, *uploadFilename, timeout)
		if err != nil {
			logger.Errorf(ctx, "Failed to upload file: %s", err)
			os.Exit(1)
		}
	case "list":
		err = listCmd.Parse(os.Args[2:])
		if err != nil {
			logger.Errorf(ctx, "Failed to parse command, please try again")
			os.Exit(1)
		}
		list, err := client.GetFileList(ctx, c, timeout)
		if err != nil {
			logger.Errorf(ctx, "Failed to get files info: %s", err)
			os.Exit(1)
		}
		for _, file := range list {
			fmt.Println(file)
		}
	case "info":
		err = infoCmd.Parse(os.Args[2:])
		if err != nil {
			logger.Errorf(ctx, "Failed to parse command, please try again")
			os.Exit(1)
		}
		if err != nil {
			logger.Errorf(ctx, "Failed to parse command, please try again")
			os.Exit(1)
		}
		if *infoFilename == "" {
			logger.Errorf(ctx, "Please provide a filename to get information.")
			os.Exit(1)
		}
		info, err := client.GetFileInfo(ctx, c, *infoFilename, timeout)
		if err != nil {
			logger.Errorf(ctx, "Failed to get file info: %s", err)
			os.Exit(1)
		}
		fmt.Printf("File size: %d, file name: %s\n", info.GetSize(), info.GetFilename())

	default:
		logger.Errorf(ctx, "Invalid command.")
		os.Exit(1)
	}
}
