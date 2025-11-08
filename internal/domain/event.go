package domain

import "github.com/bizio/abc-user-service/internal/domain/model"

type EventType string

const (
	UserCreatedEvent EventType = "UserCreated"
	UserUpdatedEvent EventType = "UserUpdated"
	UserDeletedEvent EventType = "UserDeleted"
)

type Event struct {
	Type   EventType
	UserID string
	User   *model.User
}
