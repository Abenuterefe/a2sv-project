package storage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStorage interface {
	SaveProfilePicture(file multipart.File, fileHeader *multipart.FileHeader) (string, error)
}

type localFileStorage struct {
	basePath string
}

func NewLocalFileStorage(basePath string) FileStorage {
	return &localFileStorage{basePath: basePath}
}

func (s *localFileStorage) SaveProfilePicture(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
	}
	if !allowedExts[ext] {
		return "", errors.New("invalid file type, only jpg/jpeg/png allowed")
	}

	// Ensure profile_pictures folder exists
	uploadDir := filepath.Join(s.basePath, "profile_pictures")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", errors.New("failed to create upload directory")
	}

	// Generate unique file name
	newFileName := fmt.Sprintf("%s%s", primitive.NewObjectID().Hex(), ext)
	filePath := filepath.Join(uploadDir, newFileName)

	// Create and write the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", errors.New("failed to save file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", errors.New("failed to write file")
	}

	// Return relative path so we don't leak system paths
	relativePath := filepath.ToSlash(filepath.Join("uploads", "profile_pictures", newFileName))
	return relativePath, nil
}
