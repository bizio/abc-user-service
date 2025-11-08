package service

import (
	"testing"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
	"github.com/bizio/abc-user-service/mocks"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetUserApplicationService_Do(t *testing.T) {
	userID := "user-123"

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		service := NewGetUserApplicationService(mockRepo)

		// The user that the repository is expected to return
		expectedUser, _ := model.NewUser("Test User", "test@example.com", "1999-12-31")
		expectedUser.ID = userID

		mockRepo.On("Get", userID).Return(expectedUser, nil).Once()

		res, err := service.Do(userID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, expectedUser.ToDTO(), res.User)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		service := NewGetUserApplicationService(mockRepo)

		notFoundID := "not-found-id"

		mockRepo.On("Get", notFoundID).Return(nil, domain.ErrUserNotFound).Once()

		res, err := service.Do(notFoundID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Equal(t, &v1.GetUserResponse{}, res)
		mockRepo.AssertExpectations(t)
	})
}
