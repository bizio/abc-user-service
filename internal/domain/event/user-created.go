package event

import (
	"github.com/bizio/abc-user-service/internal/domain"
	"github.com/bizio/abc-user-service/internal/domain/model"
)

func NewUserCreatedEvent(user *model.User) *domain.Event {
	return &domain.Event{Type: domain.UserCreatedEvent, UserID: user.ToDTO().ID, User: user}
}
