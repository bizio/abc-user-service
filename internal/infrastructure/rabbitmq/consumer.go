package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/bizio/abc-user-service/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	queueName    string
	exchangeName string
	channel      *amqp.Channel
}

func NewRabbitMQConsumer(queueName, exchangeName string, channel *amqp.Channel) *RabbitMQConsumer {
	return &RabbitMQConsumer{queueName, exchangeName, channel}
}

func (c *RabbitMQConsumer) Consume() (<-chan *domain.Event, error) {

	eventCh := make(chan *domain.Event)
	q, err := c.channel.QueueDeclare(c.queueName, false, false, true, false, nil)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return eventCh, err
	}
	err = c.channel.QueueBind(q.Name, "", c.exchangeName, false, nil)
	if err != nil {
		log.Printf("Failed to bind a queue: %v", err)
		return eventCh, err
	}

	msgsCh, err := c.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return eventCh, err
	}

	go func() {
		for d := range msgsCh {
			log.Printf("Received a message: %s", d.Body)
			var event domain.Event
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				continue
			}
			eventCh <- &event
		}
	}()

	log.Printf("Consumer started, waiting for messages.")
	return eventCh, nil

}
