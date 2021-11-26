package main

import (
	"context"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/Dyleme/image-coverter/pkg/handler"
	"github.com/Dyleme/image-coverter/pkg/logging"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/Dyleme/image-coverter/pkg/server"
	"github.com/Dyleme/image-coverter/pkg/service"
	"github.com/Dyleme/image-coverter/pkg/storage"
)

const imageConvertorsAmount = 5

func main() {
	logger := logging.NewLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

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

	useMinioSSL, err := strconv.ParseBool(os.Getenv("MNUSESSL"))
	if err != nil {
		logger.Fatalf("can't convert string to bool; %s", err)
	}

	stor, err := storage.NewMinioStorage(
		os.Getenv("MNHOST")+":"+os.Getenv("MNPORT"),
		os.Getenv("MNACCESSKEYID"),
		os.Getenv("MNSECRETACCESSKEY"),
		useMinioSSL,
	)

	if err != nil {
		logger.Fatalf("failed to initialize storage: %s", err)
	}

	authService := service.NewAuthSevice(authRep, &service.HashGen{}, &service.JwtGen{})
	reqService := service.NewRequestService(reqRep, stor, imageConvertorsAmount)
	downService := service.NewDownloadSerivce(downRep, stor)
	handlers := handler.New(authService, reqService, downService, logger)

	port := os.Getenv("PORT")
	srv := new(server.Server)

	ctx := logging.WithLogger(context.Background(), logger)

	if err := srv.Run(ctx, port, handlers.InitRouters()); err != nil {
		logger.Fatalf("error occurred runnging http server: %s", err.Error())
	}
}
