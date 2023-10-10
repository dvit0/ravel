package handlers

import (
	"errors"
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
	machine, found, err := h.worker.GetMachine(id)
	if err != nil {
		log.Error("failed to get machine with error :", err)
		utils.AnswerWithInternalServerError(w, err)
		return
	}
	if !found {
		log.Error("machine not found")
		utils.AnswerWithNotFoundError(w, errors.New("machine not found"))
		return
	}

	utils.AnswerWithJSON(w, machine, http.StatusOK)
}
