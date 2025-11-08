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
