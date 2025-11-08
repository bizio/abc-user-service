package local

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
)

const PathTemplate = "/user/%s/files/"

type LocalFileRepository struct {
	basePath string
}

func NewLocalFileRepository(basePath string) *LocalFileRepository {
	err := os.MkdirAll(path.Clean(basePath), os.ModePerm)
	if err != nil {
		panic(err)
	}
	return &LocalFileRepository{basePath: basePath}
}

func (s *LocalFileRepository) Upload(userID string, fileHeader *multipart.FileHeader) (string, error) {
	filePath := s.generatePath(userID, fileHeader.Filename)

	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		log.Printf("error creating directory: %s", err)
		return "", err
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("error opening file: %s", err)
		return "", err
	}
	defer file.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("error creating file: %s", err)
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("error copying file: %s", err)
		return "", err
	}

	return filePath, nil
}

func (s *LocalFileRepository) Get(userID, filename string) (*os.File, error) {
	filePath := s.generatePath(userID, filename)
	return os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
}

func (s *LocalFileRepository) List(userID string) ([]string, error) {
	list, err := os.ReadDir(s.generatePath(userID, ""))
	if err != nil {
		return nil, err
	}
	files := make([]string, len(list))
	for i, entry := range list {
		files[i] = s.generatePath(userID, entry.Name())
	}

	return files, nil
}

func (s *LocalFileRepository) Delete(userID, filename string) error {
	return os.Remove(s.generatePath(userID, filename))
}

func (s *LocalFileRepository) DeleteFiles(userID string) error {
	list, err := os.ReadDir(s.generatePath(userID, ""))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range list {
		os.Remove(s.generatePath(userID, entry.Name()))
	}

	return nil
}

func (s *LocalFileRepository) generatePath(userID, filename string) string {
	return path.Clean(s.basePath + fmt.Sprintf(PathTemplate, userID) + filename)
}
