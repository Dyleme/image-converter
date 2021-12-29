package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/Dyleme/image-coverter/internal/config"
	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/Dyleme/image-coverter/internal/rabbitmq"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/storage"
)

func main() {
	logger := logging.NewLogger(logrus.DebugLevel)

	conf, err := config.InitConfig()
	if err != nil {
		logger.Fatal("wrong config: %w", err)
	}

	ctx := logging.WithLogger(context.Background(), logger)

	db, err := repository.NewPostgresDB(conf.DB)
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err)
	}

	convRep := repository.NewConvPostgres(db)

	stor, err := storage.NewAwsStorage(conf.AwsBucketName, conf.AWS)
	if err != nil {
		logger.Fatalf("failed to initialize storage: %s", err)
	}

	convService := service.NewConvertRequest(convRep, stor)

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-c
		logger.Info("system interrupt call")
		cancel()
	}()

	err = rabbitmq.Receive(ctx, convService, conf.RabbitMQ)
	if err != nil {
		logger.Fatalf("receiving: %s", err)
	}
}
