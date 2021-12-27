package main

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/Dyleme/image-coverter/internal/config"
	"github.com/Dyleme/image-coverter/internal/handler"
	"github.com/Dyleme/image-coverter/internal/jwt"
	"github.com/Dyleme/image-coverter/internal/logging"
	"github.com/Dyleme/image-coverter/internal/rabbitmq"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/Dyleme/image-coverter/internal/server"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/storage"
)

func main() {
	logger := logging.NewLogger(logrus.InfoLevel)

	config, err := config.InitConfig()
	if err != nil {
		logger.Fatalf("wrong config: %s", err)
	}

	db, err := repository.NewPostgresDB(config.DB)
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err)
	}

	authRep := repository.NewAuthPostgres(db)
	reqRep := repository.NewReqPostgres(db)
	downRep := repository.NewDownloadPostgres(db)

	stor, err := storage.NewAwsStorage(config.AwsBucketName, config.AWS)
	fmt.Println(config.AWS, config.AwsBucketName)
	fmt.Println(config.AWS.Credentials)
	if err != nil {
		logger.Fatalf("failed to initialize storage: %s", err)
	}

	rabbitSender, err := rabbitmq.NewRabbitSender(config.RabbitMQ)
	if err != nil {
		// logger.Fatalf("failed to make connection to rabbitmq: %s", err)

	}

	jwtGen := jwt.NewJwtGen(config.JWT)

	authService := service.NewAuthSevice(authRep, &service.HashGen{}, jwtGen)
	reqService := service.NewRequestService(reqRep, stor, rabbitSender)
	downService := service.NewDownloadService(downRep, stor)

	authHandler := handler.NewAuthHandler(authService, logger)
	reqHandler := handler.NewReqHandler(reqService, logger)
	downHandler := handler.NewDownHandler(downService, logger)

	handlers := handler.New(authHandler, reqHandler, downHandler, logger)

	srv := new(server.Server)

	ctx := logging.WithLogger(context.Background(), logger)

	if err := srv.Run(ctx, config.Port, handlers.InitRouters(jwtGen)); err != nil {
		logger.Fatalf("error occurred runnging http server: %s", err.Error())
	}
}
