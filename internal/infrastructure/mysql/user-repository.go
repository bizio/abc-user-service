package mysql

import (
	"errors"

	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is the GORM model for a user
type User struct {
	gorm.Model
	ID    string `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"uniqueIndex,size:255"`
	DOB   string
	Files []*File `gorm:"foreignKey:UserID"`
}

// File is the GORM model for a file
type File struct {
	gorm.Model
	ID     string `gorm:"primaryKey"`
	UserID string `gorm:"size:255"`
	Name   string
	Path   string
	Size   int64
}

// MysqlUserRepository is the GORM implementation of the user repository
type MysqlUserRepository struct {
	db *gorm.DB
}

// NewMysqlUserRepository creates a new repository instance, runs migrations
func NewMysqlUserRepository(db *gorm.DB) *MysqlUserRepository {
	if err := db.AutoMigrate(&User{}, &File{}); err != nil {
		panic(err)
	}
	return &MysqlUserRepository{db: db}
}

// toDomainUser converts a GORM user to a domain user
func toDomainUser(u *User) *model.User {
	domainUser, _ := model.NewUser(u.Name, u.Email, u.DOB)
	domainUser.ID = u.ID
	for _, f := range u.Files {
		domainUser.AddFile(toDomainFile(f))
	}
	return domainUser
}

// fromDomainUser converts a domain user to a GORM user
func fromDomainUser(u *model.User) *User {
	files := make([]*File, 0, len(u.GetFiles()))
	for _, f := range u.GetFiles() {
		files = append(files, fromDomainFile(f))
	}
	return &User{
		ID:    u.ID,
		Name:  u.ToDTO().Name, // DTO contains the private fields
		Email: u.ToDTO().Email,
		DOB:   u.ToDTO().DOB,
		Files: files,
	}
}

// toDomainFile converts a GORM file to a domain file
func toDomainFile(f *File) *model.File {
	return &model.File{
		ID:     f.ID,
		UserID: f.UserID,
		Name:   f.Name,
		Path:   f.Path,
		Size:   f.Size,
	}
}

// fromDomainFile converts a domain file to a GORM file
func fromDomainFile(f *model.File) *File {
	return &File{
		ID:     f.ID,
		UserID: f.UserID,
		Name:   f.Name,
		Path:   f.Path,
		Size:   f.Size,
	}
}

func (r *MysqlUserRepository) Create(user *model.User) (string, error) {
	user.ID = uuid.NewString()
	persistenceUser := fromDomainUser(user)

	result := r.db.Create(persistenceUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return "", domain.ErrUserAlreadyExists
		}
		return "", result.Error
	}
	return persistenceUser.ID, nil
}

func (r *MysqlUserRepository) Get(id string) (*model.User, error) {
	var user User
	result := r.db.Preload("Files").First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}
	return toDomainUser(&user), nil
}

func (r *MysqlUserRepository) GetByEmail(email string) (*model.User, error) {
	var user User
	result := r.db.Preload("Files").First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}
	return toDomainUser(&user), nil
}

func (r *MysqlUserRepository) List() ([]*model.User, error) {
	var users []User
	result := r.db.Preload("Files").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	var domainUsers []*model.User
	for _, u := range users {
		domainUsers = append(domainUsers, toDomainUser(&u))
	}
	return domainUsers, nil
}

func (r *MysqlUserRepository) Update(id string, user *model.User) error {
	// Find the existing record to preserve fields like created_at
	var existingUser User
	if err := r.db.First(&existingUser, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return err
	}

	updatedPersistenceUser := fromDomainUser(user)

	existingUser.Name = updatedPersistenceUser.Name
	existingUser.Email = updatedPersistenceUser.Email
	existingUser.DOB = updatedPersistenceUser.DOB
	existingUser.Files = updatedPersistenceUser.Files

	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&existingUser).Error
}

func (r *MysqlUserRepository) Delete(id string) error {
	// soft delete user and associated files
	return r.db.Select("Files").Delete(&User{ID: id}).Error
}

func (r *MysqlUserRepository) GetFiles(userID string) ([]*model.File, error) {
	var files []File
	result := r.db.Where("user_id = ?", userID).Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}
	var domainFiles []*model.File
	for _, f := range files {
		domainFiles = append(domainFiles, toDomainFile(&f))
	}
	return domainFiles, nil
}

func (r *MysqlUserRepository) GetFile(userID, fileID string) (*model.File, error) {
	var file File
	result := r.db.First(&file, "user_id = ? AND id = ?", userID, fileID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound // Or a more specific file not found error
		}
		return nil, result.Error
	}
	return toDomainFile(&file), nil
}

func (r *MysqlUserRepository) DeleteFiles(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&File{}).Error
}
