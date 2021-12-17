package main

import (
	"context"
	"os"

	_ "github.com/lib/pq"

	"github.com/Dyleme/image-coverter/internal/handler"
	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/Dyleme/image-coverter/internal/rabbitmq"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/Dyleme/image-coverter/internal/server"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/storage"
)

func main() {
	logger := logging.NewLogger()

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

	authRep := repository.NewAuthPostgres(db)
	reqRep := repository.NewReqPostgres(db)
	downRep := repository.NewDownloadPostgres(db)

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

	rabbitSender, err := rabbitmq.NewRabbitSender(rabbitConfig)
	if err != nil {
		logger.Fatalf("failed to make connection to rabbitmq: %s", err)
	}

	authService := service.NewAuthSevice(authRep, &service.HashGen{}, &service.JwtGen{})
	reqService := service.NewRequestService(reqRep, stor, rabbitSender)
	downService := service.NewDownloadService(downRep, stor)

	handlers := handler.New(authService, reqService, downService, logger)

	port := os.Getenv("PORT")
	srv := new(server.Server)

	ctx := logging.WithLogger(context.Background(), logger)

	if err := srv.Run(ctx, port, handlers.InitRouters()); err != nil {
		logger.Fatalf("error occurred runnging http server: %s", err.Error())
	}
}
