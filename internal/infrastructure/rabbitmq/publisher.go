package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/bizio/abc-user-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	exchangeName string
	channel      *amqp.Channel
}

func NewRabbitMQPublisher(exchangeName string, channel *amqp.Channel) *RabbitMQPublisher {
	return &RabbitMQPublisher{exchangeName, channel}
}

func (p *RabbitMQPublisher) Publish(event *domain.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.channel.ExchangeDeclare(p.exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare exchange: %v", err)
		return err
	}

	encodedEvent, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to encode event: %v", err)
		return err
	}
	err = p.channel.PublishWithContext(ctx, p.exchangeName, "", false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        encodedEvent,
		})
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
		return err
	}

	log.Printf("Published event of type %s for user ID %s", event.Type, event.UserID)
	return nil
}
