package main

import (
	"context"
	"os"

	_ "github.com/lib/pq"

	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/rabbitmq"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/storage"
)

type emptySender struct{}

func (r *emptySender) ProcessImage(data *model.ConversionData) {
}

func main() {
	logger := logging.NewLogger()
	ctx := logging.WithLogger(context.Background(), logger)

	db, err := repository.NewPostgresDB(&repository.DBConfig{
		UserName: os.Getenv("DBUSERNAME"),
		Password: os.Getenv("DBPASSWORD"),
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("DBSSLMODE"),
	})
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err)
	}

	reqRep := repository.NewReqPostgres(db)
	stor, err := storage.NewAwsStorage()

	if err != nil {
		logger.Fatalf("failed to initialize storage: %s", err)
	}

	rabbitConfig := rabbitmq.Config{
		User:     os.Getenv("RBUSER"),
		Password: os.Getenv("RBPASSWORD"),
		Host:     os.Getenv("RBHOST"),
		Port:     os.Getenv("RBPORT"),
	}

	reqService := service.NewRequestService(reqRep, stor, &emptySender{})

	rabbitmq.Receive(ctx, reqService, reqService, rabbitConfig)
}
