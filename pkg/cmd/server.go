package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/bizio/abc-user-service/internal/infrastructure/rabbitmq"
	"github.com/bizio/abc-user-service/pkg/protocol/rest"
	env "github.com/caarlos0/env/v11"
	amqp "github.com/rabbitmq/amqp091-go"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	HTTPPort            string `env:"HTTP_PORT"`
	DatastoreDBHost     string `env:"DB_HOST"`
	DatastoreDBPort     string `env:"DB_PORT"`
	DatastoreDBUser     string `env:"DB_USER"`
	DatastoreDBPassword string `env:"DB_PASSWORD"`
	DatastoreDBName     string `env:"DB_NAME"`
	QueueUser           string `env:"QUEUE_USER"`
	QueuePassword       string `env:"QUEUE_PASSWORD"`
	QueueHost           string `env:"QUEUE_HOST"`
	QueuePort           string `env:"QUEUE_PORT"`
}

// RunServer runs HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("failed to parse environment variables: %s", err)
		return err
	}

	param := "charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBPort,
		cfg.DatastoreDBName,
		param)

	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to connect to mysql: %s", err)
		panic(err)
	}

	queueConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.QueueUser,
		cfg.QueuePassword,
		cfg.QueueHost,
		cfg.QueuePort,
	)
	amqpConn, err := amqp.Dial(queueConnectionString)
	if err != nil {
		log.Printf("failed to connect to RabbitMQ: %s", err)
		panic(err)
	}
	defer amqpConn.Close()

	channel, err := amqpConn.Channel()
	if err != nil {
		log.Printf("failed to open a channel: %s", err)
		panic(err)
	}
	defer channel.Close()

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP server: '%s'", cfg.HTTPPort)
	}

	rabbitmqConsumer := rabbitmq.NewRabbitMQConsumer("user_events_queue", "user_events", channel)
	eventCh, err := rabbitmqConsumer.Consume()
	if err != nil {
		log.Printf("Failed to start RabbitMQ consumer: %v", err)
		panic(err)
	}

	go func() {
		for event := range eventCh {
			log.Printf("Processing event: %v", event)
			// Add event processing logic here
		}
	}()

	fmt.Printf("Starting HTTP/REST gateway on port %s...\n", cfg.HTTPPort)
	return rest.RunServer(ctx, cfg.HTTPPort, db, channel)
}
