package service

import (
	"log"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/event"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
)

func NewUpdateUserApplicationService(repository domain.UserRepository, publisher domain.EventPublisher) *UpdateUserApplicationService {
	return &UpdateUserApplicationService{repository, publisher}
}

type UpdateUserApplicationService struct {
	repository domain.UserRepository
	publisher  domain.EventPublisher
}

func (s *UpdateUserApplicationService) Do(req *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {

	user, err := s.repository.Get(req.ID)
	if err != nil {
		return &v1.UpdateUserResponse{}, err
	}

	if req.Name != "" {
		user.SetName(req.Name)
	}

	if req.Email != "" {
		err = user.SetEmail(req.Email)
		if err != nil {
			return &v1.UpdateUserResponse{}, err
		}
	}

	if req.DOB != "" {
		err = user.SetDob(req.DOB)
		if err != nil {
			return &v1.UpdateUserResponse{}, err
		}
	}

	err = s.repository.Update(req.ID, user)
	if err != nil {
		return &v1.UpdateUserResponse{}, err
	}

	go func() {
		err = s.publisher.Publish(event.NewUserUpdatedEvent(user))
		if err != nil {
			log.Printf("Failed to publish user created event: %v", err)
		}
	}()

	return &v1.UpdateUserResponse{User: user.ToDTO()}, nil

}
