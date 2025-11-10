package service

import (
	"log"

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

	go func() {
		err = s.publisher.Publish(event.NewUserDeletedEvent(id))
		if err != nil {
			log.Printf("Failed to publish user created event: %v", err)
		}
	}()
	return nil

}
