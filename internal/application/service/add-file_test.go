package service

import (
	"errors"
	"mime/multipart"
	"testing"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
	"github.com/bizio/abc-user-service/mocks"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddFileApplicationService_Do(t *testing.T) {
	userID := "user-123"
	maxSize := int64(1024)

	// Create a base user for tests
	user, _ := model.NewUser("Test User", "test@example.com", "1990-01-01")
	user.ID = userID

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewAddFileApplicationService(mockUserRepo, mockFileRepo, maxSize)

		fileHeader := &multipart.FileHeader{
			Filename: "test.jpg",
			Size:     512,
		}
		req := &v1.UploadFileRequest{UserID: userID, File: fileHeader}
		filePath := "/uploads/test.jpg"

		// Need a fresh copy of the user for the mock return
		userCopy, _ := model.NewUser("Test User", "test@example.com", "1990-01-01")
		userCopy.ID = userID

		mockUserRepo.On("Get", userID).Return(userCopy, nil).Once()
		mockFileRepo.On("Upload", userID, fileHeader).Return(filePath, nil).Once()
		mockUserRepo.On("Update", userID, mock.Anything).Return(nil).Once()

		res, err := service.Do(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, fileHeader.Filename, res.File.Name)
		assert.Equal(t, filePath, res.File.Path)
		assert.NotEmpty(t, res.File.ID)
		mockUserRepo.AssertExpectations(t)
		mockFileRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewAddFileApplicationService(mockUserRepo, mockFileRepo, maxSize)

		req := &v1.UploadFileRequest{UserID: "not-found"}

		mockUserRepo.On("Get", "not-found").Return(nil, domain.ErrUserNotFound).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, res)
		mockFileRepo.AssertNotCalled(t, "Upload", mock.Anything, mock.Anything)
		mockUserRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("File Too Large", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewAddFileApplicationService(mockUserRepo, mockFileRepo, maxSize)

		fileHeader := &multipart.FileHeader{Size: maxSize + 1}
		req := &v1.UploadFileRequest{UserID: userID, File: fileHeader}

		userCopy, _ := model.NewUser("Test User", "test@example.com", "1990-01-01")
		userCopy.ID = userID

		mockUserRepo.On("Get", userID).Return(userCopy, nil).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, model.ErrFileTooLarge)
		assert.Nil(t, res)
		mockFileRepo.AssertNotCalled(t, "Upload", mock.Anything, mock.Anything)
	})

	t.Run("Storage Upload Fails", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewAddFileApplicationService(mockUserRepo, mockFileRepo, maxSize)

		fileHeader := &multipart.FileHeader{Size: 512}
		req := &v1.UploadFileRequest{UserID: userID, File: fileHeader}
		uploadErr := errors.New("s3 upload failed")

		userCopy, _ := model.NewUser("Test User", "test@example.com", "1990-01-01")
		userCopy.ID = userID

		mockUserRepo.On("Get", userID).Return(userCopy, nil).Once()
		mockFileRepo.On("Upload", userID, fileHeader).Return("", uploadErr).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, uploadErr)
		assert.Nil(t, res)
		mockUserRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("User Update Fails", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		mockFileRepo := new(mocks.FileRepository)
		service := NewAddFileApplicationService(mockUserRepo, mockFileRepo, maxSize)

		fileHeader := &multipart.FileHeader{Size: 512}
		req := &v1.UploadFileRequest{UserID: userID, File: fileHeader}
		updateErr := errors.New("db update failed")

		userCopy, _ := model.NewUser("Test User", "test@example.com", "1990-01-01")
		userCopy.ID = userID

		mockUserRepo.On("Get", userID).Return(userCopy, nil).Once()
		mockFileRepo.On("Upload", userID, fileHeader).Return("/path", nil).Once()
		mockUserRepo.On("Update", userID, mock.Anything).Return(updateErr).Once()

		res, err := service.Do(req)

		assert.ErrorIs(t, err, updateErr)
		assert.Nil(t, res)
		mockUserRepo.AssertExpectations(t)
		mockFileRepo.AssertExpectations(t)
	})
}
