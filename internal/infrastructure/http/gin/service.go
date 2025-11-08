package http

import (
	"net/http"

	applicationService "github.com/bizio/abc-user-service/internal/application/service"
	"github.com/bizio/abc-user-service/internal/domain"
	v1 "github.com/bizio/abc-user-service/pkg/api/v1"
	"github.com/gin-gonic/gin"
)

type HttpError struct {
	Error string `json:"error"`
}

type GinHttpService struct {
	listService        *applicationService.ListUsersApplicationService
	getService         *applicationService.GetUserApplicationService
	createService      *applicationService.CreateUserApplicationService
	updateService      *applicationService.UpdateUserApplicationService
	deleteService      *applicationService.DeleteUserApplicationService
	getFilesSerivce    *applicationService.GetFilesApplicationService
	addFileService     *applicationService.AddFileApplicationService
	deleteFilesService *applicationService.DeleteFilesApplicationService
	maxFileSize        int64
}

func NewGinHttpService(
	listService *applicationService.ListUsersApplicationService,
	getService *applicationService.GetUserApplicationService,
	createService *applicationService.CreateUserApplicationService,
	updateService *applicationService.UpdateUserApplicationService,
	deleteService *applicationService.DeleteUserApplicationService,
	getFilesService *applicationService.GetFilesApplicationService,
	addFileService *applicationService.AddFileApplicationService,
	deleteApplicationService *applicationService.DeleteFilesApplicationService,
	maxFileSize int64,
) *GinHttpService {
	return &GinHttpService{
		listService,
		getService,
		createService,
		updateService,
		deleteService,
		getFilesService,
		addFileService,
		deleteApplicationService,
		maxFileSize,
	}

}

func (s *GinHttpService) GetRouter() http.Handler {
	router := gin.Default()

	router.MaxMultipartMemory = s.maxFileSize
	v1Users := router.Group("/v1/users")
	v1Users.GET("", s.List)
	v1Users.GET("/:id", s.Get)
	v1Users.POST("", s.Create)
	v1Users.PUT("/:id", s.Update)
	v1Users.DELETE("/:id", s.Delete)
	v1Users.GET("/:id/files", s.GetFiles)
	v1Users.POST("/:id/files", s.UploadFile)
	v1Users.DELETE("/:id/files", s.DeleteFiles)

	return router
}

// List list all users
//	@Summary		List all users
//	@Description	List all users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	v1.ListUsersResponse
//	@Failure		500		{object}	HttpError
//	@Router			/users [GET]
func (s *GinHttpService) List(c *gin.Context) {
	users, err := s.listService.Do()
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// Get get a user by ID
//	@Summary		Get a user by ID
//	@Description	Get a single user by its ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	v1.User
//	@Failure		404	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/users/{id} [GET]
func (s *GinHttpService) Get(c *gin.Context) {
	req := v1.GetUserRequest{}
	if err := c.BindUri(&req); err != nil {
		handleError(c, err)
		return
	}

	user, err := s.getService.Do(req.ID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Create create a new user
//	@Summary		Create a new user
//	@Description	Create a new user with the provided information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		v1.CreateUserRequest	true	"User to create"
//	@Success		201		{object}	v1.CreateUserResponse
//	@Failure		400		{object}	HttpError
//	@Failure		500		{object}	HttpError
//	@Router			/users [POST]
func (s *GinHttpService) Create(c *gin.Context) {
	req := &v1.CreateUserRequest{}

	if err := c.BindJSON(req); err != nil {
		handleError(c, err)
		return
	}

	res, err := s.createService.Do(req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Update update a user
//	@Summary		Update a user
//	@Description	Update an existing user's information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"User ID"
//	@Param			user	body		v1.UpdateUserRequest	true	"User data to update"
//	@Success		201		{object}	v1.UpdateUserResponse
//	@Failure		400		{object}	HttpError
//	@Failure		404		{object}	HttpError
//	@Failure		500		{object}	HttpError
//	@Router			/users/{id} [PUT]
func (s *GinHttpService) Update(c *gin.Context) {
	req := &v1.UpdateUserRequest{}

	if err := c.BindJSON(req); err != nil {
		handleError(c, err)
		return
	}

	res, err := s.updateService.Do(req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)

}

// Delete delete a user
//	@Summary		Delete a user
//	@Description	Delete a user by its ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		204	{object}	nil
//	@Failure		404	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/users/{id} [DELETE]
func (s *GinHttpService) Delete(c *gin.Context) {
	req := v1.DeleteUserRequest{}

	if err := c.BindUri(&req); err != nil {
		handleError(c, err)
		return
	}

	err := s.deleteService.Do(req.ID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetFiles get user's files
//	@Summary		Get user's files
//	@Description	Get a list of files for a specific user
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	v1.GetFilesResponse
//	@Failure		404	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/users/{id}/files [GET]
func (s *GinHttpService) GetFiles(c *gin.Context) {
	req := v1.GetFilesRequest{}

	if err := c.BindUri(&req); err != nil {
		handleError(c, err)
		return
	}

	files, err := s.getFilesSerivce.Do(req.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, files)

}

// UploadFile upload a file for a user
//	@Summary		Upload a file
//	@Description	Upload a file for a specific user
//	@Tags			files
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string	true	"User ID"
//	@Param			file	formData	file	true	"File to upload"
//	@Success		201		{object}	v1.UploadFileResponse
//	@Failure		400		{object}	HttpError
//	@Failure		500		{object}	HttpError
//	@Router			/users/{id}/files [POST]
func (s *GinHttpService) UploadFile(c *gin.Context) {
	req := &v1.UploadFileRequest{}
	if err := c.Bind(&req); err != nil {
		handleError(c, err)
		return
	}

	res, err := s.addFileService.Do(req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// DeleteFiles delete all files for a user
//	@Summary		Delete all files for a user
//	@Description	Delete all files associated with a specific user
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		204	{object}	nil
//	@Failure		404	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/users/{id}/files [DELETE]
func (s *GinHttpService) DeleteFiles(c *gin.Context) {
	req := v1.DeleteFilesRequest{}

	if err := c.BindUri(&req); err != nil {
		handleError(c, err)
		return
	}

	err := s.deleteFilesService.Do(req.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)

}

func handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
