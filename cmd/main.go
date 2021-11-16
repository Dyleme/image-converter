package main

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/Dyleme/image-coverter"
	"github.com/Dyleme/image-coverter/pkg/handler"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/Dyleme/image-coverter/pkg/service"
	"github.com/Dyleme/image-coverter/pkg/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
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
		log.Fatalf("failed to initialize db: %s", err)
	}

	authRep := repository.NewAuthPostgres(db)
	reqRep := repository.NewReqPostgres(db)
	downRep := repository.NewDownloadPostgres(db)

	useMinioSSL, err := strconv.ParseBool(os.Getenv("MNUSESSL"))
	if err != nil {
		log.Fatalf("can't convert string to bool; %s", err)
	}

	stor, err := storage.NewMinioStorage(
		os.Getenv("MNHOST")+":"+os.Getenv("MNPORT"),
		os.Getenv("MNACCESSKEYID"),
		os.Getenv("MNSECRETACCESSKEY"),
		useMinioSSL,
	)

	if err != nil {
		log.Fatalf("failed to initialize storage: %s", err)
	}

	authService := service.NewAuthSevice(authRep)
	reqService := service.NewRequestService(reqRep, stor)
	downService := service.NewDownloadSerivce(downRep, stor)
	handlers := handler.NewServer(authService, reqService, downService)

	port := os.Getenv("PORT")
	srv := new(image.Server)

	if err := srv.Run(port, handlers.InitRouters()); err != nil {
		log.Fatalf("error occurred runnging http server: %s", err.Error())
	}
}
