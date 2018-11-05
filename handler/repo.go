package handler

import "net/http"

type repository struct{}

func (h *repository) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("TODO"))
}
