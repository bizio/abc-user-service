package domain

//go:generate mockery --name EventConsumer --output ../../mocks --outpkg mocks
type EventConsumer interface {
	Consume() (<-chan *Event, error)
}
