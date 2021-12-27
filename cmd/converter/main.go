package main

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/Dyleme/image-coverter/internal/config"
	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/Dyleme/image-coverter/internal/rabbitmq"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/storage"
)

type emptySender struct{}

func (r *emptySender) ProcessImage(ctx context.Context, data *rabbitmq.ConversionData) {
}

func main() {
	logger := logging.NewLogger(logrus.DebugLevel)

	config, err := config.InitConfig()
	if err != nil {
		logger.Fatal("wrong config: %w", err)
	}

	ctx := logging.WithLogger(context.Background(), logger)

	db, err := repository.NewPostgresDB(config.DB)
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err)
	}

	reqRep := repository.NewReqPostgres(db)

	stor, err := storage.NewAwsStorage(config.AwsBucketName, config.AWS)
	if err != nil {
		logger.Fatalf("failed to initialize storage: %s", err)
	}

	reqService := service.NewRequestService(reqRep, stor, &emptySender{})

	rabbitmq.Receive(ctx, reqService, config.RabbitMQ)
}
