package service

import (
	"errors"
	"testing"

	"github.com/bizio/abc-user-service/internal/domain/model"
	"github.com/bizio/abc-user-service/mocks"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestListUsersApplicationService_Do(t *testing.T) {

	t.Run("Success with multiple users", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		service := NewListUsersApplicationService(mockRepo)

		user1, _ := model.NewUser("User One", "one@example.com", "1991-01-01")
		user1.ID = "user-1"
		user2, _ := model.NewUser("User Two", "two@example.com", "1992-02-02")
		user2.ID = "user-2"

		expectedUsers := []*model.User{user1, user2}

		mockRepo.On("List").Return(expectedUsers, nil).Once()

		res, err := service.Do()

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, int32(2), res.Count)
		assert.Len(t, res.Users, 2)
		assert.Equal(t, user1.ToDTO(), res.Users[0])
		assert.Equal(t, user2.ToDTO(), res.Users[1])
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success with no users", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		service := NewListUsersApplicationService(mockRepo)

		expectedUsers := []*model.User{} // Empty slice

		mockRepo.On("List").Return(expectedUsers, nil).Once()

		res, err := service.Do()

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, int32(0), res.Count)
		assert.Len(t, res.Users, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		service := NewListUsersApplicationService(mockRepo)

		repoErr := errors.New("database connection lost")

		mockRepo.On("List").Return(nil, repoErr).Once()

		res, err := service.Do()

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, &v1.ListUsersResponse{}, res)
		mockRepo.AssertExpectations(t)
	})
}
