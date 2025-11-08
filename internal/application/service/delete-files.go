package service

import (
	"github.com/bizio/abc-user-service/internal/domain"
)

func NewDeleteFilesApplicationService(repository domain.UserRepository, storage domain.FileRepository) *DeleteFilesApplicationService {
	return &DeleteFilesApplicationService{repository, storage}
}

type DeleteFilesApplicationService struct {
	repository domain.UserRepository
	storage    domain.FileRepository
}

func (s *DeleteFilesApplicationService) Do(userID string) error {
	user, err := s.repository.Get(userID)
	if err != nil {
		return err
	}

	err = s.storage.DeleteFiles(user.ID)
	if err != nil {
		return err
	}

	err = s.repository.DeleteFiles(user.ID)
	if err != nil {
		return err
	}

	return nil

}
