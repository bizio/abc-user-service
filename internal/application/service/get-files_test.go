package service

import (
	"testing"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
	"github.com/bizio/abc-user-service/mocks"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetFilesApplicationService_Do(t *testing.T) {
	userID := "user-123"

	t.Run("Success with files", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		service := NewGetFilesApplicationService(mockUserRepo)

		// Create a user with some files
		user, _ := model.NewUser("Test User", "test@example.com", "1990-01-01")
		user.ID = userID
		file1 := &model.File{ID: "file-1", Name: "photo.jpg"}
		file2 := &model.File{ID: "file-2", Name: "resume.pdf"}
		user.AddFile(file1)
		user.AddFile(file2)

		mockUserRepo.On("Get", userID).Return(user, nil).Once()

		res, err := service.Do(userID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Files, 2)
		assert.Equal(t, file1.ToDTO(), res.Files[0])
		assert.Equal(t, file2.ToDTO(), res.Files[1])
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Success with no files", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		service := NewGetFilesApplicationService(mockUserRepo)

		// Create a user with no files
		user, _ := model.NewUser("Test User", "test@example.com", "1990-01-01")
		user.ID = userID

		mockUserRepo.On("Get", userID).Return(user, nil).Once()

		res, err := service.Do(userID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Files, 0)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockUserRepo := new(mocks.UserRepository)
		service := NewGetFilesApplicationService(mockUserRepo)

		mockUserRepo.On("Get", userID).Return(nil, domain.ErrUserNotFound).Once()

		res, err := service.Do(userID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Equal(t, &v1.GetFilesResponse{}, res)
		mockUserRepo.AssertExpectations(t)
	})
}
