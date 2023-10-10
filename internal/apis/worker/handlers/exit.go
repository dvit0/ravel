package handlers

import (
	"net/http"
	"os"

	"github.com/valyentdev/ravel/internal/utils"
)

func (h *Handler) ExitWorker(w http.ResponseWriter, r *http.Request) {

	utils.AnswerWithJSON(w, map[string]string{"message": "worker is shutting down..."}, http.StatusOK)
	defer os.Exit(0)
}
