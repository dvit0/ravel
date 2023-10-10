package main

import (
	"os"
	"os/signal"
	"syscall"

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	worker := worker.NewWorker()

	server := api.CreateWorkerHTTPServer(worker)

	go func() {
		<-c
		log.Info("Shutting down...")
		server.Close()
	}()

	if err := server.ListenAndServe(); err != nil {
		log.Error(err)
	}

}
