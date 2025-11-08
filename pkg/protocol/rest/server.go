package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	service "github.com/bizio/abc-user-service/internal/application/service"
	infraHttp "github.com/bizio/abc-user-service/internal/infrastructure/http/gin"
	"github.com/bizio/abc-user-service/internal/infrastructure/mysql"
	"github.com/bizio/abc-user-service/internal/infrastructure/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"

	"github.com/bizio/abc-user-service/internal/infrastructure/storage/local"
)

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, httpPort string, db *gorm.DB, channel *amqp.Channel) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var maxFileSize int64 = 2 << 20 // 2 MB
	localFileRepository := local.NewLocalFileRepository(os.TempDir())
	mysqlRepository := mysql.NewMysqlUserRepository(db)
	rabbitmqPublisher := rabbitmq.NewRabbitMQPublisher("user_events", channel)

	listApplicationService := service.NewListUsersApplicationService(mysqlRepository)
	getApplicationService := service.NewGetUserApplicationService(mysqlRepository)
	createApplicationService := service.NewCreateUserApplicationService(mysqlRepository, rabbitmqPublisher)
	updateApplicationService := service.NewUpdateUserApplicationService(mysqlRepository, rabbitmqPublisher)
	deleteApplicationService := service.NewDeleteUserApplicationService(mysqlRepository, localFileRepository, rabbitmqPublisher)

	getFilesApplicationService := service.NewGetFilesApplicationService(mysqlRepository)
	addFileApplicationService := service.NewAddFileApplicationService(mysqlRepository, localFileRepository, int64(maxFileSize))
	deleteFilesApplicationService := service.NewDeleteFilesApplicationService(mysqlRepository, localFileRepository)

	httpService := infraHttp.NewGinHttpService(
		listApplicationService, getApplicationService, createApplicationService, updateApplicationService,
		deleteApplicationService, getFilesApplicationService, addFileApplicationService, deleteFilesApplicationService,
		maxFileSize,
	)

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: httpService.GetRouter(),
	}
	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server forced to shutdown: %v", err)
		}
	}()

	log.Println("starting HTTP/REST gateway on port " + httpPort)
	return srv.ListenAndServe()
}
