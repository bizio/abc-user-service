package domain

import (
	"mime/multipart"
	"os"
)

//go:generate mockery --name FileRepository --output ../../mocks --outpkg mocks
type FileRepository interface {
	Upload(userID string, file *multipart.FileHeader) (string, error)
	Get(userID, filename string) (*os.File, error)
	List(userID string) ([]string, error)
	Delete(userID, filename string) error
	DeleteFiles(userID string) error
}
