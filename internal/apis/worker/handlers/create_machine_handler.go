package handlers

import (
	"context"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/utils"
	"github.com/valyentdev/ravel/pkg/types"
)

func (h *Handler) CreateMachineHandler(w http.ResponseWriter, r *http.Request) {
	var machineSpec types.RavelMachineSpec

	if err := utils.DecodeJSON(w, r, &machineSpec); err != nil {
		log.Error("failed to decode json", "err", err)
		utils.AnswerWithBadRequestError(w, err)
		return
	}

	id, err := h.worker.CreateMachine(context.Background(), machineSpec)

	if err != nil {
		log.Error("failed to create machine", "err", err)
		utils.AnswerWithInternalServerError(w, err)
		return
	}

	utils.AnswerWithJSON(w, map[string]string{"id": id}, http.StatusCreated)
}
