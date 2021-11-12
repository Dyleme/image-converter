package main

import (
	"log"
	"os"
	"strconv"

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

	repos := repository.NewRepository(db)

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

	services := service.NewService(repos, stor)
	handlers := handler.NewServer(services)

	port := os.Getenv("PORT")
	srv := new(image.Server)

	if err := srv.Run(port, handlers.InitRouters()); err != nil {
		log.Fatalf("error occurred runnging http server: %s", err.Error())
	}
}
