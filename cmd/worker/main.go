package main

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	api "github.com/valyentdev/ravel/internal/apis/worker"
	"github.com/valyentdev/ravel/internal/worker"
)

func main() {
	log.SetPrefix("[Worker]")
	log.Info("Ravel worker starting...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	worker := worker.NewWorker()
	api.StartWorkerApi(worker)
}
