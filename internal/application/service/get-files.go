package service

import (
	"github.com/bizio/abc-user-service/internal/domain"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
)

func NewGetFilesApplicationService(repository domain.UserRepository) *GetFilesApplicationService {
	return &GetFilesApplicationService{repository}
}

type GetFilesApplicationService struct {
	repository domain.UserRepository
}

func (s *GetFilesApplicationService) Do(userID string) (*v1.GetFilesResponse, error) {
	user, err := s.repository.Get(userID)
	if err != nil {
		return &v1.GetFilesResponse{}, err
	}

	files := make([]*v1.File, 0, len(user.GetFiles()))
	for _, file := range user.GetFiles() {
		files = append(files, file.ToDTO())
	}

	return &v1.GetFilesResponse{Files: files}, nil

}
