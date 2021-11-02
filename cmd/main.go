package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/Dyleme/image-coverter"
	"github.com/Dyleme/image-coverter/pkg/controller"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/Dyleme/image-coverter/pkg/service"
	"github.com/Dyleme/image-coverter/pkg/storage"
)

func main() {
	err := godotenv.Load("C:\\Users\\Aliaksei.Dziauho\\program\\image-coverter\\.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	db, err := repository.NewPostgresDB(&repository.DBConfig{
		UserName: os.Getenv("DBUSERNAME"),
		Password: os.Getenv("DBPASSWORD"),
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("DBPORT"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("DBSSLMODE"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err)
	}

	repos := repository.NewRepository(db)
	stor := storage.NewStorage()
	services := service.NewService(repos, stor)
	handlers := controller.NewController(services)

	port := os.Getenv("PORT")
	srv := new(image.Server)

	if err := srv.Run(port, handlers.InitRouters()); err != nil {
		log.Fatalf("error occurred runnging http server: %s", err.Error())
	}
}
