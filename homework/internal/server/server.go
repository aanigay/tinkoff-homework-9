package server

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework/internal/file"
	"homework/internal/logger"
	"io"

	pb "homework/pkg/files"
)

// Интерфейс для работы с файлами на сервере.
type fileService interface {
	saveFile(file file.File) error
	getFileList() ([]string, error)
	getFileInfo(string) (*file.Info, error)
}

type Server struct {
	pb.UnimplementedFileServiceServer
	fs fileService
}

func NewServer(fs fileService) *Server {
	return &Server{
		fs: fs,
	}
}

func (s *Server) UploadFile(stream pb.FileService_UploadFileServer) error {
	ctx := stream.Context()

	// Получение файла от клиента.
	var fileData []byte
	filename := ""
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf(ctx, "failed to read file chunk: %s", err)
			return err
		}

		filename = chunk.Filename
		fileData = append(fileData, chunk.ChunkData...)
	}
	f := file.File{
		Name:     filename,
		Contents: fileData,
	}
	// Сохранение файла на сервере.
	err := s.fs.saveFile(f)
	if err != nil {
		logger.Errorf(ctx, "failed to save file: %s", err)
		return err
	}

	logger.Infof(ctx, "successfully uploaded file: %s", filename)

	// Отправка ответа клиенту.
	return stream.SendAndClose(&pb.UploadResponse{Success: true, Message: "File uploaded successfully"})
}

func (s *Server) GetFileList(ctx context.Context, empty *emptypb.Empty) (*pb.FileList, error) {
	// Получение списка файлов на сервере.
	files, err := s.fs.getFileList()
	if err != nil {
		logger.Errorf(ctx, "failed to get files list: %s", err)
		return nil, err
	}

	return &pb.FileList{Files: files}, nil
}

func (s *Server) GetFileInfo(ctx context.Context, request *pb.FileInfoRequest) (*pb.FileInfo, error) {
	// Получение информации о файле на сервере.
	fileInfo, err := s.fs.getFileInfo(request.Filename)
	if err != nil {
		logger.Errorf(ctx, "failed to get file info: %s", err)
		return nil, err
	}
	info := &pb.FileInfo{
		Filename: fileInfo.Name,
		Size:     fileInfo.Size,
	}
	return info, nil
}
