package service

import (
	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/event"
)

func NewDeleteUserApplicationService(repository domain.UserRepository, storage domain.FileRepository, publisher domain.EventPublisher) *DeleteUserApplicationService {
	return &DeleteUserApplicationService{repository, storage, publisher}
}

type DeleteUserApplicationService struct {
	repository domain.UserRepository
	storage    domain.FileRepository
	publisher  domain.EventPublisher
}

func (s *DeleteUserApplicationService) Do(id string) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}

	err = s.storage.DeleteFiles(id)
	if err != nil {
		return err
	}

	s.publisher.Publish(event.NewUserDeletedEvent(id))

	return nil

}
