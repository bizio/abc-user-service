package service

import (
	"errors"
	"testing"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
	"github.com/bizio/abc-user-service/mocks"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserApplicationService_Do(t *testing.T) {
	t.Run("DOB Validation Error", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewCreateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.CreateUserRequest{
			Name:  "Jane Doe",
			Email: "john.doe@example.com",
			DOB:   "2025-01-33",
		}
		res, err := service.Do(req)

		assert.ErrorIs(t, err, model.ErrInvalidDob)
		assert.Equal(t, &v1.CreateUserResponse{}, res, "should return empty response on error")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("Email Validation Error", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewCreateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.CreateUserRequest{
			Name:  "Jane Doe",
			Email: "invalid-email",
			DOB:   "2000-01-01",
		}
		res, err := service.Do(req)

		assert.ErrorIs(t, err, model.ErrInvalidEmailAddress)
		assert.Equal(t, &v1.CreateUserResponse{}, res, "should return empty response on error")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("User Already Exists Error", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewCreateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.CreateUserRequest{
			Name:  "Jane Doe",
			Email: "duplicate@example.com",
			DOB:   "2000-01-01",
		}
		user, _ := model.NewUser(req.Name, req.Email, req.DOB)

		mockRepo.On("GetByEmail", req.Email).Return(user, nil).Once()
		mockRepo.On("Create", user).Return("", nil).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
		assert.Equal(t, &v1.CreateUserResponse{}, res, "should return empty response on error")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewCreateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.CreateUserRequest{
			Name:  "John Doe",
			Email: "john.doe@example.com",
			DOB:   "2000-01-01",
		}

		repoErr := errors.New("unexpected database error")

		mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return("", repoErr).Once()
		mockRepo.On("GetByEmail", req.Email).Return(nil, domain.ErrUserNotFound).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, repoErr)
		assert.Equal(t, &v1.CreateUserResponse{}, res)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.UserRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewCreateUserApplicationService(mockRepo, mockEventPublisher)

		req := &v1.CreateUserRequest{
			Name:  "John Doe",
			Email: "john.doe@example.com",
			DOB:   "2000-01-01",
		}

		expectedUserID := "new-user-id-123"

		// Assert that the Create method is called with a user matching the request data.
		user, _ := model.NewUser(req.Name, req.Email, req.DOB)
		mockRepo.On("Create", user).Return(expectedUserID, nil).Once()
		mockRepo.On("GetByEmail", req.Email).Return(nil, domain.ErrUserNotFound).Once()
		mockEventPublisher.On("Publish", mock.Anything).Return(nil).Once()
		res, err := service.Do(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, expectedUserID, res.ID)
		mockRepo.AssertExpectations(t)
	})

}
