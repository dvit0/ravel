package handlers

import (
	"context"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/utils"
)

func (h *Handler) DeleteMachineHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	err := h.worker.DeleteMachine(context.Background(), id)

	if err != nil {
		log.Error("failed to delete machine with error :", err)
		utils.AnswerWithInternalServerError(w, err)
		return
	}

	utils.AnswerWithJSON(w, map[string]string{"message": "Machine deleted"}, http.StatusCreated)
}
