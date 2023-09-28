package handlers

import (
	"context"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/utils"
)

func (h *Handler) StartMachineHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	if id == "" {
		log.Error("failed to get id from url")
		utils.AnswerWithBadRequestError(w, nil)
		return
	}

	err := h.worker.StartMachine(context.Background(), id)

	if err != nil {
		log.Error("failed to create machine with error : ", err)
		utils.AnswerWithInternalServerError(w, err)
		return
	}

	utils.AnswerWithJSON(w, map[string]string{"message": "Machine started"}, http.StatusOK)
}
