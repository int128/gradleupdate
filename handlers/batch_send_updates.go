package handlers

import (
	"net/http"

	"github.com/int128/gradleupdate/usecases/interfaces"
	"google.golang.org/appengine/log"
)

type BatchSendUpdates struct {
	BatchSendUpdates usecases.BatchSendUpdates
}

func (h *BatchSendUpdates) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := h.BatchSendUpdates.Do(ctx); err != nil {
		log.Errorf(ctx, "error while sending updates: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
