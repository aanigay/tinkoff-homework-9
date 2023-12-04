package server

import (
	"fmt"
	"homework/internal/file"
	"io/fs"
	"os"
	"path/filepath"
)

type FileHandler struct {
	fileStoragePath string
}

func NewHandler(path string) *FileHandler {
	return &FileHandler{
		fileStoragePath: path,
	}
}

// Сохранение файла на сервере.
func (h *FileHandler) saveFile(file file.File) error {
	filePath := filepath.Join(h.fileStoragePath, file.Name)
	err := os.WriteFile(filePath, file.Contents, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}
	return nil
}

// Получение списка файлов.
func (h *FileHandler) getFileList() ([]string, error) {
	entries, err := os.ReadDir(h.fileStoragePath)
	if err != nil {
		return nil, err
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read info list: %v", err)
	}

	var fileList []string
	for _, info := range infos {
		fileList = append(fileList, info.Name())
	}

	return fileList, nil
}

// Получение информации об указанном файле.
func (h *FileHandler) getFileInfo(filename string) (*file.Info, error) {
	filePath := filepath.Join(h.fileStoragePath, filename)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	return &file.Info{
		Name: filename,
		Size: fileInfo.Size(),
	}, nil
}
