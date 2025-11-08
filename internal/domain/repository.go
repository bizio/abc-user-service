package domain

import (
	"errors"

	"github.com/bizio/abc-user-service/internal/domain/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

//go:generate mockery --name UserRepository --output ../../mocks --outpkg mocks
type UserRepository interface {
	Create(user *model.User) (string, error)
	Get(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	List() ([]*model.User, error)
	Update(id string, user *model.User) error
	Delete(id string) error
	GetFiles(userID string) ([]*model.File, error)
	GetFile(userID, fileID string) (*model.File, error)
	DeleteFiles(userID string) error
}
