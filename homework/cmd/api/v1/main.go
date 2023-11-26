package main

import (
	"github.com/gorilla/mux"
	"homework/internal/config"
	"homework/internal/handlers"
	"homework/internal/repository"
	"homework/internal/usecase/impl"
	"log"
	"net/http"
)

func main() {
	// инициализация Hanlder, Service
	// запуск http сервера
	c, err := config.Read("cfg.yaml")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := mux.NewRouter()
	repo := repository.New()
	deviceUC := impl.New(repo)
	handler := handlers.NewHandler(deviceUC)
	handler.RegisterHandlers(router)

	log.Fatal(http.ListenAndServe(config.ServerAddress(*c), router))
}
