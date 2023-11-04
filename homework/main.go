package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	c, err := config.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := mux.NewRouter()
	repo := repository.New()
	deviceUC := impl.New(repo)
	handler := handlers.NewHandler(deviceUC)
	handler.RegisterHandlers(router)

	log.Fatal(http.ListenAndServe(c.ServerAddress(), router))
}
