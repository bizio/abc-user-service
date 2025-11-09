package domain

//go:generate mockery --name EventPublisher --output ../../mocks --outpkg mocks
type EventPublisher interface {
	Publish(event *Event) error
}
