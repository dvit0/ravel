package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/valyentdev/ravel/internal/utils"
	"github.com/valyentdev/ravel/pkg/types"
)

type ListMachinesPayload struct {
	Machines []types.RavelMachine `json:"machines"`
}

func (h *Handler) ListMachinesHandler(w http.ResponseWriter, r *http.Request) {
	machines, err := h.worker.ListMachines()

	if err != nil {
		log.Error("failed to list machines with error : ", err)
		utils.AnswerWithInternalServerError(w, err)
		return
	}

	utils.AnswerWithJSON(w, ListMachinesPayload{Machines: machines}, http.StatusOK)
}
