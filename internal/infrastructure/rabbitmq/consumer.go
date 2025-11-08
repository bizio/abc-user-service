package rabbitmq

import (
	"log"

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

func (c *RabbitMQConsumer) Consume() error {

	q, err := c.channel.QueueDeclare(c.queueName, false, false, true, false, nil)

	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}
	err = c.channel.QueueBind(q.Name, "", c.exchangeName, false, nil)
	if err != nil {
		log.Printf("Failed to bind a queue: %v", err)
		return err
	}

	msgsCh, err := c.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return err
	}

	go func() {
		for d := range msgsCh {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("Consumer started, waiting for messages.")
	return nil

}
