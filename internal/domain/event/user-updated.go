package event

import (
	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
)

func NewUserUpdatedEvent(user *model.User) *domain.Event {
	return &domain.Event{Type: domain.UserUpdatedEvent, UserID: user.ToDTO().ID, User: user}
}
