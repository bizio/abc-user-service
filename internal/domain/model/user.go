package model

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"time"

	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
)

const MinimumAge = 18

var (
	ErrInvalidEmailAddress     = errors.New("invalid email address")
	ErrInvalidDob              = errors.New("invalid date of birth")
	ErrMinAgeRequirementNotMet = fmt.Errorf("user must be at least %d years old", MinimumAge)
)

type User struct {
	ID    string
	name  string
	email string
	dob   string
	files []*File
}

func NewUser(name string, email string, dob string) (*User, error) {

	user := &User{name: name}
	err := user.SetDob(dob)
	if err != nil {
		return nil, err
	}

	err = user.SetEmail(email)
	if err != nil {
		return nil, err
	}
	user.files = make([]*File, 0)
	return user, nil
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetDob(dob string) error {
	parsedDob, err := time.Parse(time.DateOnly, dob)
	if err != nil {
		log.Printf("error parsing date: %s", err)
		return ErrInvalidDob
	}

	if time.Now().Year()-parsedDob.Year() < MinimumAge {
		return ErrMinAgeRequirementNotMet
	}

	u.dob = parsedDob.Format(time.DateOnly)
	return nil
}

func (u *User) SetEmail(email string) error {
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		log.Printf("error parsing email address: %s", err)
		return ErrInvalidEmailAddress
	}

	u.email = parsedEmail.Address

	return nil
}

func (u *User) AddFile(file *File) {
	u.files = append(u.files, file)
}

func (u *User) DeleteFile(fileID string) {
	for i, file := range u.files {
		if file.ID == fileID {
			u.files = append(u.files[:i], u.files[i+1:]...)
			return
		}
	}
}

func (u *User) DeleteFiles() {
	u.files = make([]*File, 0)
}

func (u *User) GetFiles() []*File {
	return u.files
}

func (u *User) ToDTO() *v1.User {
	files := make([]*v1.File, 0, len(u.files))
	for _, f := range u.files {
		files = append(files, f.ToDTO())
	}
	return &v1.User{
		ID:    u.ID,
		Name:  u.name,
		Email: u.email,
		DOB:   u.dob,
		Files: files,
	}
}
