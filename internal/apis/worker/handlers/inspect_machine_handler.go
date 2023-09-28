package handlers

import (
	"context"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/utils"
)

func (h *Handler) GetMachineHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")
	if id == "" {
		log.Error("failed to get id from url")
		utils.AnswerWithInternalServerError(w, nil)
		return
	}
	machine, err := h.worker.GetMachine(context.Background(), id)

	if err != nil {
		log.Error("failed to get machine with error :", err)
		utils.AnswerWithNotFoundError(w, err)
		return
	}

	utils.AnswerWithJSON(w, machine, http.StatusOK)
}
