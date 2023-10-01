package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jeronimobarea/go-cqrs/internal/feed"
)

type handler struct {
	feedsSvc feed.Service
}

func (h handler) createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var cmd feed.CreateFeedCmd
	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.feedsSvc.Create(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}
