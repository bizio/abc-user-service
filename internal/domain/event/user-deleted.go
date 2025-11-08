package event

import "github.com/bizio/abc-user-service/internal/domain"

func NewUserDeletedEvent(userID string) *domain.Event {
	return &domain.Event{UserID: userID, Type: domain.UserDeletedEvent}
}
