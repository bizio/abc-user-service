package model

import (
	"errors"

	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
)

var ErrFileTooLarge = errors.New("file is too large")

type File struct {
	ID     string
	UserID string
	Name   string
	Path   string
	Size   int64
}

func (f *File) ToDTO() *v1.File {
	return &v1.File{
		ID:     f.ID,
		UserID: f.UserID,
		Name:   f.Name,
		Path:   f.Path,
		Size:   f.Size,
	}
}
