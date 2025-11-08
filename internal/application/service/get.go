package service

import (
	"github.com/bizio/abc-user-service/internal/domain"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
)

func NewGetUserApplicationService(repository domain.UserRepository) *GetUserApplicationService {
	return &GetUserApplicationService{repository}
}

type GetUserApplicationService struct {
	repository domain.UserRepository
}

func (s *GetUserApplicationService) Do(id string) (*v1.GetUserResponse, error) {
	user, err := s.repository.Get(id)
	if err != nil {
		return &v1.GetUserResponse{}, err
	}

	return &v1.GetUserResponse{User: user.ToDTO()}, nil

}
