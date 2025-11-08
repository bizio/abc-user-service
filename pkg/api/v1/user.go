package v1

import "mime/multipart"

// DTOs
type User struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
	DOB   string  `json:"dob"`
	Files []*File `json:"files"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	DOB   string `json:"dob" binding:"required"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type UpdateUserRequest struct {
	ID    string `json:"id" uri:"id" binding:"required"`
	Name  string `json:"name" binding:"omitempty"`
	Email string `json:"email" binding:"omitempty,email"`
	DOB   string `json:"dob" binding:"omitempty"`
}

type UpdateUserResponse struct {
	User *User `json:"user"`
}

type ListUsersResponse struct {
	Users []*User `json:"users"`
	Count int32   `json:"count"`
}

type GetUserRequest struct {
	ID string `json:"id" uri:"id" binding:"required"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

type DeleteUserRequest struct {
	ID string `json:"id" uri:"id" binding:"required"`
}

type File struct {
	ID     string `json:"id"`
	UserID string `json:"userID"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Size   int64  `json:"size"`
}

type GetFilesRequest struct {
	UserID string `json:"id" uri:"id" binding:"required"`
}

type GetFilesResponse struct {
	Files []*File `json:"files"`
}

type UploadFileRequest struct {
	UserID string                `form:"id" uri:"id" binding:"required"`
	File   *multipart.FileHeader `form:"file" binding:"required"`
}

type UploadFileResponse struct {
	File *File `json:"file"`
}

type DeleteFilesRequest struct {
	UserID string `json:"id" uri:"id" binding:"required"`
}
