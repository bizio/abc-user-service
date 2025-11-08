package service

import (
	"errors"
	"testing"
	"time"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
	"github.com/bizio/abc-user-service/mocks"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateUserApplicationService_Do(t *testing.T) {
	userID := "user-123"

	// Create a valid existing user for tests.
	existingUser, _ := model.NewUser("Old Name", "old.email@example.com", "1990-01-01")
	existingUser.ID = userID

	t.Run("Success", func(t *testing.T) {
		// Need a fresh user copy for each sub-test to avoid data races
		userCopy, _ := model.NewUser("Old Name", "old.email@example.com", "1990-01-01")
		userCopy.ID = userID

		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewUpdateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.UpdateUserRequest{
			ID:    userID,
			Name:  "New Name",
			Email: "new.email@example.com",
			DOB:   "1991-02-02",
		}

		mockRepo.On("Get", userID).Return(userCopy, nil).Once()
		mockRepo.On("Update", userID, mock.Anything).Return(nil).Once()
		mockEventPublisher.On("Publish", mock.Anything).Return(nil).Once()

		res, err := service.Do(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.Name, res.User.Name)
		assert.Equal(t, req.Email, res.User.Email)
		assert.Equal(t, req.DOB, res.User.DOB)
		mockRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewUpdateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.UpdateUserRequest{ID: "not-found-id"}

		mockRepo.On("Get", "not-found-id").Return(nil, domain.ErrUserNotFound).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Equal(t, &v1.UpdateUserResponse{}, res)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("Invalid Email on Update", func(t *testing.T) {
		userCopy, _ := model.NewUser("Old Name", "old.email@example.com", "1990-01-01")
		userCopy.ID = userID

		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewUpdateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.UpdateUserRequest{
			ID:    userID,
			Email: "this-is-not-an-email",
		}

		mockRepo.On("Get", userID).Return(userCopy, nil).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, model.ErrInvalidEmailAddress)
		assert.Equal(t, &v1.UpdateUserResponse{}, res)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("Invalid DOB on Update", func(t *testing.T) {
		userCopy, _ := model.NewUser("Old Name", "old.email@example.com", "1990-01-01")
		userCopy.ID = userID

		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewUpdateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.UpdateUserRequest{
			ID:  userID,
			DOB: "not-a-date",
		}

		mockRepo.On("Get", userID).Return(userCopy, nil).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, model.ErrInvalidDob)
		assert.Equal(t, &v1.UpdateUserResponse{}, res)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("Age requirment not met on Update", func(t *testing.T) {
		underAgeDate := time.Date(time.Now().Year()-17, 1, 1, 0, 0, 0, 0, time.Local)

		userCopy, _ := model.NewUser("Old Name", "old.email@example.com", "1990-01-01")
		userCopy.ID = userID

		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewUpdateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.UpdateUserRequest{
			ID:  userID,
			DOB: underAgeDate.Format(time.DateOnly),
		}

		mockRepo.On("Get", userID).Return(userCopy, nil).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, model.ErrMinAgeRequirementNotMet)
		assert.Equal(t, &v1.UpdateUserResponse{}, res)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("Repository Update Fails", func(t *testing.T) {
		userCopy, _ := model.NewUser("Old Name", "old.email@example.com", "1990-01-01")
		userCopy.ID = userID

		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewUpdateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.UpdateUserRequest{
			ID:   userID,
			Name: "A New Name",
		}

		repoErr := errors.New("db-update-failed")

		mockRepo.On("Get", userID).Return(userCopy, nil).Once()
		mockRepo.On("Update", userID, userCopy).Return(repoErr).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, &v1.UpdateUserResponse{}, res)
		mockRepo.AssertExpectations(t)
	})
}
