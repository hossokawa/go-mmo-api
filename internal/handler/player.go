package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hossokawa/go-nethttp-example/internal/player"
)

type PlayerHandler struct {
	service *player.PlayerService
}

func NewPlayerHandler(service *player.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: service}
}

type IDParsingError struct {
	value string
}

func NewIDParsingError(value string) error {
	return &IDParsingError{value: value}
}

func (e *IDParsingError) Error() string {
	return fmt.Sprintf("parsing id: invalid value '%v'", e.value)
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params player.CreatePlayerParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if params.Username == "" || params.Class == "" {
		http.Error(w, "Username and/or class cannot be empty", http.StatusBadRequest)
		return
	}

	p, err := h.service.CreatePlayer(context.Background(), params.Username, params.Class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *PlayerHandler) GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/`json")

	players, err := h.service.GetAllPlayers(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(players)
}

func (h *PlayerHandler) GetPlayerByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, NewIDParsingError(idStr).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	p, err := h.service.GetPlayerByID(context.Background(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}
