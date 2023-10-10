package api

import (
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/valyentdev/ravel/internal/apis/worker/handlers"
	"github.com/valyentdev/ravel/internal/worker"
)

const (
	PORT = ":3000"
)

func CreateWorkerHTTPServer(worker *worker.Worker) *http.Server {

	h := handlers.NewHandler(worker)
	mux := flow.New()
	mux.Use(authMiddleware)

	mux.HandleFunc("/api/v1/exit", h.ExitWorker, "POST")

	mux.HandleFunc("/api/v1/machines", h.CreateMachineHandler, "POST")
	mux.HandleFunc("/api/v1/machines", h.ListMachinesHandler, "GET")

	mux.HandleFunc("/api/v1/machines/:id", h.GetMachineHandler, "GET")
	mux.HandleFunc("/api/v1/machines/:id/start", h.StartMachineHandler, "POST")
	mux.HandleFunc("/api/v1/machines/:id/stop", h.StopMachineHandler, "POST")
	mux.HandleFunc("/api/v1/machines/:id", h.DeleteMachineHandler, "DELETE")

	server := &http.Server{
		Addr:    PORT,
		Handler: http.HandlerFunc(mux.ServeHTTP),
	}
	return server
}
