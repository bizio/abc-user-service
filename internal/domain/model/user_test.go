package model

import (
	"testing"
	"time"

	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	underAgeDate := time.Date(time.Now().Year()-17, 1, 1, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name          string
		userName      string
		email         string
		dob           string
		expectedEmail string
		expectedErr   error
	}{
		{
			name:          "Valid User",
			userName:      "John Doe",
			email:         "john.doe@example.com",
			dob:           "2000-01-01",
			expectedEmail: "john.doe@example.com",
			expectedErr:   nil,
		},
		{
			name:        "Invalid Email",
			userName:    "Jane Doe",
			email:       "invalid-email",
			dob:         "2000-01-01",
			expectedErr: ErrInvalidEmailAddress,
		},
		{
			name:        "Invalid DOB",
			userName:    "Jane Doe",
			email:       "jane.doe@example.com",
			dob:         "not-a-date",
			expectedErr: ErrInvalidDob,
		},
		{
			name:        "Age requirement not met",
			userName:    "Jane Doe",
			email:       "jane.doe@example.com",
			dob:         underAgeDate.Format(time.DateOnly),
			expectedErr: ErrMinAgeRequirementNotMet,
		},
		{
			name:          "Email with name part",
			userName:      "John Doe",
			email:         "John Doe <john.doe@example.com>",
			dob:           "2000-01-01",
			expectedEmail: "john.doe@example.com",
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.userName, tt.email, tt.dob)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userName, user.name)
				assert.Equal(t, tt.expectedEmail, user.email)
				assert.Equal(t, tt.dob, user.dob)
			}
		})
	}
}

func TestUser_SetName(t *testing.T) {
	user := &User{name: "Old Name"}
	user.SetName("New Name")
	assert.Equal(t, "New Name", user.name)
}

func TestUser_SetEmail(t *testing.T) {
	user := &User{email: "jane.doe@example.com"}
	err := user.SetEmail("john.doe@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "john.doe@example.com", user.email)
}
func TestUser_SetDob(t *testing.T) {
	user := &User{dob: "2000-01-01"}
	err := user.SetDob("1990-12-31")
	assert.NoError(t, err)
	assert.Equal(t, "1990-12-31", user.dob)
}
func TestUser_AddFile(t *testing.T) {
	user := &User{ID: "user-123", files: []*File{}}
	file := &File{ID: "file-456", UserID: "user-123", Name: "example.txt", Path: "/tmp/user/user-123/files/example.txt", Size: 256}
	user.AddFile(file)

	assert.Len(t, user.files, 1)
	assert.Equal(t, file, user.files[0])
}

func TestUser_DeleteFile(t *testing.T) {
	user := &User{ID: "user-123", files: []*File{
		{ID: "file-123", UserID: "user-123", Name: "file1.txt", Path: "/tmp/user/user-123/files/file1.txt", Size: 128},
		{ID: "file-456", UserID: "user-123", Name: "file2.txt", Path: "/tmp/user/user-123/files/file2.txt", Size: 256},
	}}
	user.DeleteFile("file-123")
	assert.Len(t, user.files, 1)
	assert.Equal(t, "file-456", user.files[0].ID)

	user.DeleteFile("non-existent-file")
	assert.Len(t, user.files, 1)
}

func TestUser_DeleteFiles(t *testing.T) {
	user := &User{ID: "user-123", files: []*File{
		{ID: "file-123", UserID: "user-123", Name: "file1.txt", Path: "/tmp/user/user-123/files/file1.txt", Size: 128},
		{ID: "file-456", UserID: "user-123", Name: "file2.txt", Path: "/tmp/user/user-123/files/file2.txt", Size: 256},
	}}
	user.DeleteFiles()
	assert.Len(t, user.files, 0)
}

func TestUser_GetFiles(t *testing.T) {
	user := &User{ID: "user-123", files: []*File{
		{ID: "file-123", UserID: "user-123", Name: "file1.txt", Path: "/tmp/user/user-123/files/file1.txt", Size: 128},
		{ID: "file-456", UserID: "user-123", Name: "file2.txt", Path: "/tmp/user/user-123/files/file2.txt", Size: 256},
	}}
	files := user.GetFiles()
	assert.Len(t, files, 2)
	assert.Equal(t, "file-123", files[0].ID)
	assert.Equal(t, "file-456", files[1].ID)
}
func TestUser_ToDTO(t *testing.T) {
	user := &User{
		ID:    "user-123",
		name:  "Test User",
		email: "test@example.com",
		dob:   "1995-05-10",
		files: []*File{
			{ID: "123-456", UserID: "user-123", Name: "example.txt", Path: "/tmp/user/user-123/files/example.txt", Size: 128},
		},
	}

	expectedDto := &v1.User{
		ID:    "user-123",
		Name:  "Test User",
		Email: "test@example.com",
		DOB:   "1995-05-10",
		Files: []*v1.File{
			{ID: "123-456", UserID: "user-123", Name: "example.txt", Path: "/tmp/user/user-123/files/example.txt", Size: 128},
		},
	}

	dto := user.ToDTO()

	assert.Equal(t, expectedDto, dto)
}
