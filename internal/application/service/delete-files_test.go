package service

import (
	"errors"
	"testing"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
	"github.com/bizio/abc-user-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteFilesApplicationService_Do(t *testing.T) {
	userID := "user-123"

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewDeleteFilesApplicationService(mockUserRepo, mockFileRepo)

		user, _ := model.NewUser("Test", "test@test.com", "1990-01-01")
		user.ID = userID

		mockUserRepo.On("Get", userID).Return(user, nil).Once()
		mockFileRepo.On("DeleteFiles", userID).Return(nil).Once()
		mockUserRepo.On("DeleteFiles", userID).Return(nil).Once()

		err := service.Do(userID)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockFileRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewDeleteFilesApplicationService(mockUserRepo, mockFileRepo)

		mockUserRepo.On("Get", userID).Return(nil, domain.ErrUserNotFound).Once()

		err := service.Do(userID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockFileRepo.AssertNotCalled(t, "DeleteFiles", mock.Anything)
		mockUserRepo.AssertNotCalled(t, "DeleteFiles", mock.Anything)
		mockUserRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("Storage Deletion Fails", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewDeleteFilesApplicationService(mockUserRepo, mockFileRepo)

		user, _ := model.NewUser("Test", "test@test.com", "1990-01-01")
		user.ID = userID
		storageErr := errors.New("storage error")

		mockUserRepo.On("Get", userID).Return(user, nil).Once()
		mockFileRepo.On("DeleteFiles", userID).Return(storageErr).Once()

		err := service.Do(userID)

		assert.ErrorIs(t, err, storageErr)
		mockUserRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("User Update Fails", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewDeleteFilesApplicationService(mockUserRepo, mockFileRepo)

		user, _ := model.NewUser("Test", "test@test.com", "1990-01-01")
		user.ID = userID
		updateErr := errors.New("db update error")

		mockUserRepo.On("Get", userID).Return(user, nil).Once()
		mockFileRepo.On("DeleteFiles", userID).Return(nil).Once()
		mockUserRepo.On("DeleteFiles", userID).Return(updateErr).Once()

		err := service.Do(userID)

		assert.ErrorIs(t, err, updateErr)
		mockUserRepo.AssertExpectations(t)
		mockFileRepo.AssertExpectations(t)
	})
}
