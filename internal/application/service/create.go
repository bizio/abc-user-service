package service

import (
	"log"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/event"
	"github.com/bizio/abc-user-service/internal/domain/model"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
)

func NewCreateUserApplicationService(repository domain.UserRepository, publisher domain.EventPublisher) *CreateUserApplicationService {
	return &CreateUserApplicationService{repository, publisher}
}

type CreateUserApplicationService struct {
	repository domain.UserRepository
	publisher  domain.EventPublisher
}

func (s *CreateUserApplicationService) Do(req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {

	user, err := model.NewUser(req.Name, req.Email, req.DOB)
	if err != nil {
		return &v1.CreateUserResponse{}, err
	}

	// check if user with same email already exists
	existingUser, err := s.repository.GetByEmail(req.Email)
	if existingUser != nil && err == nil {
		return &v1.CreateUserResponse{}, domain.ErrUserAlreadyExists
	}

	id, err := s.repository.Create(user)
	if err != nil {
		return &v1.CreateUserResponse{}, err
	}

	event, err := event.NewUserCreatedEvent(user)
	if err != nil {
		log.Printf("failed to publish user created event: %v", err)
	} else {
		s.publisher.Publish(event)
	}

	return &v1.CreateUserResponse{ID: id}, nil

}
