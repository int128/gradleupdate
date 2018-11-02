package handler

import "net/http"

// Home handles a request for index.
type Home struct{}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
