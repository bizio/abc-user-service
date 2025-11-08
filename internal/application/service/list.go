package service

import (
	"github.com/bizio/abc-user-service/internal/domain"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
)

func NewListUsersApplicationService(repository domain.UserRepository) *ListUsersApplicationService {
	return &ListUsersApplicationService{repository}
}

type ListUsersApplicationService struct {
	repository domain.UserRepository
}

func (s *ListUsersApplicationService) Do() (*v1.ListUsersResponse, error) {
	users, err := s.repository.List()
	if err != nil {
		return &v1.ListUsersResponse{}, err
	}

	userDTOs := make([]*v1.User, len(users))
	for i, user := range users {
		userDTOs[i] = user.ToDTO()
	}

	return &v1.ListUsersResponse{
		Users: userDTOs,
		Count: int32(len(userDTOs)),
	}, nil

}
