package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/hossokawa/go-nethttp-example/internal/api"
	"github.com/hossokawa/go-nethttp-example/internal/player"
)

type PlayerHandler struct {
	service *player.PlayerService
}

func NewPlayerHandler(service *player.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: service}
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params player.CreatePlayerParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		api.WriteJSONError(w, http.StatusBadRequest, "Error decoding request body into CreatePlayerParams struct")
		return
	}
	defer r.Body.Close()

	if params.Username == "" || params.Class == "" {
		api.WriteJSONError(w, http.StatusBadRequest, "Username and/or class cannot be empty")
		return
	}

	p, err := h.service.CreatePlayer(context.Background(), params.Username, params.Class)
	if err != nil {
		api.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *PlayerHandler) GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/`json")

	username := r.URL.Query().Get("username")

	if username != "" {
		p, err := h.service.GetPlayerByUsername(context.Background(), username)
		if err != nil {
			var notFoundErr *player.NotFoundErr
			if errors.As(err, &notFoundErr) {
				api.WriteJSONError(w, http.StatusNotFound, notFoundErr.Error())
				return
			}
			api.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(w).Encode(p)
	} else {
		players, err := h.service.GetAllPlayers(context.Background())
		if err != nil {
			api.WriteJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(w).Encode(players)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PlayerHandler) GetPlayerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.WriteJSONError(w, http.StatusInternalServerError, api.NewIDParsingError(idStr).Error())
		return
	}

	p, err := h.service.GetPlayerByID(context.Background(), int32(id))
	if err != nil {
		var notFoundErr *player.NotFoundErr
		if errors.As(err, &notFoundErr) {
			api.WriteJSONError(w, http.StatusNotFound, notFoundErr.Error())
			return
		}
		api.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}
