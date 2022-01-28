package main

import (
	"context"

	_ "github.com/lib/pq"

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
	conf, err := config.InitConfig()
	logger := logging.NewLogger(conf.LogLevel)

	if err != nil {
		logger.Fatalf("wrong config: %s", err)
	}

	db, err := repository.NewPostgresDB(conf.DB)
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err)
	}

	authRep := repository.NewAuthPostgres(db)
	reqRep := repository.NewReqPostgres(db)
	downRep := repository.NewDownloadPostgres(db)

	stor, err := storage.NewAwsStorage(conf.AwsBucketName, conf.AWS)
	if err != nil {
		logger.Fatalf("failed to initialize storage: %s", err)
	}

	rabbitSender, err := rabbitmq.NewRabbitSender(conf.RabbitMQ)
	if err != nil {
		logger.Fatalf("failed to make connection to rabbitmq: %s", err)
	}

	jwtGen := jwt.NewJwtGen(conf.JWT)

	authService := service.NewAuth(authRep, &service.HashGen{}, jwtGen)
	reqService := service.NewRequest(reqRep, stor, rabbitSender)
	downService := service.NewDownload(downRep, stor)

	authHandler := handler.NewAuth(authService, logger)
	reqHandler := handler.NewRequest(reqService, logger)
	downHandler := handler.NewDownload(downService, logger)

	handlers := handler.New(authHandler, reqHandler, downHandler, logger)

	srv := new(server.Server)

	ctx := logging.WithLogger(context.Background(), logger)

	if err := srv.Run(ctx, conf.Port, handlers.InitRouters(jwtGen)); err != nil {
		logger.Fatalf("error occurred runnging http server: %s", err.Error())
	}
}
