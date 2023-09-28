package api

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/apis/worker/handlers"
	"github.com/valyentdev/ravel/internal/worker"
)

const (
	PORT = ":3000"
)

func StartWorkerAPI(worker *worker.Worker) {
	h := handlers.NewHandler(worker)
	mux := flow.New()
	mux.Use(authMiddleware)

	mux.HandleFunc("/api/v1/machines", h.CreateMachineHandler, "POST")
	mux.HandleFunc("/api/v1/machines", h.ListMachinesHandler, "GET")
	mux.HandleFunc("/api/v1/machines/:id", h.GetMachineHandler, "GET")
	mux.HandleFunc("/api/v1/machines/:id/start", h.StartMachineHandler, "POST")
	mux.HandleFunc("/api/v1/machines/:id/stop", h.StopMachineHandler, "POST")
	mux.HandleFunc("/api/v1/machines/:id", h.DeleteMachineHandler, "DELETE")

	log.Info(fmt.Sprintf("Starting API on port %s", PORT))
	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}
