package handlers

import (
	"net/http"

	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"go.uber.org/dig"
)

type BatchSendUpdates struct {
	dig.In
	BatchSendUpdates usecases.BatchSendUpdates
	Logger           gateways.Logger
}

func (h *BatchSendUpdates) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := h.BatchSendUpdates.Do(ctx); err != nil {
		h.Logger.Errorf(ctx, "error while sending updates: %s", err)
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
}
