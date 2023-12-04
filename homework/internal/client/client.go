package client

import (
	"context"
	"fmt"
	"homework/internal/logger"
	"io"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "homework/pkg/files"
)

const (
	DefaultChunk    = 1024 * 1024 // 1 MB
	fileStoragePath = "client_files"
)

// UploadFile Загрузка файла на сервер.
func UploadFile(ctx context.Context, client pb.FileServiceClient, filename string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	stream, err := client.UploadFile(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to upload file: %s", err)
		return err
	}
	defer stream.CloseSend()

	file, err := os.Open(filepath.Join(fileStoragePath, filename))
	if err != nil {
		logger.Errorf(ctx, "failed to read file: %s", err)
		return err
	}
	defer file.Close()

	firstRead := true
	buffer := make([]byte, DefaultChunk)
	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			// Костыль на случай если файл пустой.
			if firstRead {
				chunk := &pb.FileChunk{
					Filename:  filename,
					ChunkData: make([]byte, 0),
				}
				if err := stream.Send(chunk); err != nil {
					logger.Errorf(ctx, "failed to send file chunk to server: %s", err)
					return err
				}
			}
			break
		}
		firstRead = false
		if err != nil {
			logger.Errorf(ctx, "error occurred while reading file: %s", err)
			return err
		}

		chunk := &pb.FileChunk{
			Filename:  filename,
			ChunkData: buffer[:bytesRead],
		}

		if err := stream.Send(chunk); err != nil {
			logger.Errorf(ctx, "failed to send file chunk to server: %s", err)
			return err
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf(ctx, "failed to get response: %s", err)
		return err
	}

	fmt.Printf("Upload response: %v\n", response)
	return nil
}

// GetFileList Получение списка файлов на сервере.
func GetFileList(ctx context.Context, client pb.FileServiceClient, timeout time.Duration) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	response, err := client.GetFileList(ctx, &emptypb.Empty{})
	if err != nil {
		logger.Errorf(ctx, "failed to get files list: %s", err)
		return nil, err
	}

	return response.Files, nil
}

// GetFileInfo Получение информации об указанном файле.
func GetFileInfo(ctx context.Context, client pb.FileServiceClient, filename string, timeout time.Duration) (*pb.FileInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	response, err := client.GetFileInfo(ctx, &pb.FileInfoRequest{Filename: filename})
	if err != nil {
		logger.Errorf(ctx, "failed to get file info: %s", err)
		return nil, err
	}

	return response, nil
}
