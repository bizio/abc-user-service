package service

import (
	"log"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/google/uuid"
)

func NewAddFileApplicationService(
	repository domain.UserRepository,
	storage domain.FileRepository,
	maxFileSize int64) *AddFileApplicationService {
	return &AddFileApplicationService{repository, storage, maxFileSize}
}

type AddFileApplicationService struct {
	repository  domain.UserRepository
	storage     domain.FileRepository
	maxFileSize int64
}

func (s *AddFileApplicationService) Do(req *v1.UploadFileRequest) (*v1.UploadFileResponse, error) {
	user, err := s.repository.Get(req.UserID)
	if err != nil {
		return nil, err
	}

	if req.File.Size > s.maxFileSize {
		return nil, model.ErrFileTooLarge
	}

	filepath, err := s.storage.Upload(req.UserID, req.File)
	if err != nil {
		log.Printf("error uploading file: %s", err)
		return nil, err
	}

	newFile := &model.File{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		Name:   req.File.Filename,
		Path:   filepath,
		Size:   req.File.Size,
	}
	user.AddFile(newFile)

	err = s.repository.Update(user.ID, user)
	if err != nil {
		return nil, err
	}

	return &v1.UploadFileResponse{File: user.GetFiles()[len(user.GetFiles())-1].ToDTO()}, nil

}
