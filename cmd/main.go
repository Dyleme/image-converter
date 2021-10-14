package main

import (
	"fmt"

	"github.com/Dyleme/image-coverter"
	"github.com/Dyleme/image-coverter/pkg/handler"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/Dyleme/image-coverter/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	fmt.Println("main")
	srv := new(image.Server)
	srv.Run("8080", handlers.InitRouters())
}
