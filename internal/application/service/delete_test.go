package service

import (
	"errors"
	"testing"

	"github.com/bizio/abc-user-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUserApplicationService_Do(t *testing.T) {
	userID := "user-to-delete"

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewDeleteUserApplicationService(mockUserRepo, mockFileRepo, mockEventPublisher)

		mockUserRepo.On("Delete", userID).Return(nil).Once()
		mockFileRepo.On("DeleteFiles", userID).Return(nil).Once()
		mockEventPublisher.On("Publish", mock.Anything).Return(nil).Once()

		err := service.Do(userID)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockFileRepo.AssertExpectations(t)
	})

	t.Run("User repository error", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewDeleteUserApplicationService(mockUserRepo, mockFileRepo, mockEventPublisher)

		repoErr := errors.New("user not found in db")
		mockUserRepo.On("Delete", userID).Return(repoErr).Once()

		err := service.Do(userID)

		assert.ErrorIs(t, err, repoErr)
		mockUserRepo.AssertExpectations(t)
		// Ensure file repo is not called if user repo fails
		mockFileRepo.AssertNotCalled(t, "DeleteFiles", userID)
	})

	t.Run("File repository error", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		mockEventPublisher := new(mocks.EventPublisher)
		service := NewDeleteUserApplicationService(mockUserRepo, mockFileRepo, mockEventPublisher)

		storageErr := errors.New("s3 bucket error")

		mockUserRepo.On("Delete", userID).Return(nil).Once()
		mockFileRepo.On("DeleteFiles", userID).Return(storageErr).Once()

		err := service.Do(userID)

		assert.ErrorIs(t, err, storageErr)
		mockUserRepo.AssertExpectations(t)
		mockFileRepo.AssertExpectations(t)
	})
}
