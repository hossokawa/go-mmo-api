package routes

import (
	"net/http"

	"github.com/hossokawa/go-nethttp-example/internal/handler"
	"github.com/hossokawa/go-nethttp-example/internal/player"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(router *http.ServeMux, db *pgx.Conn) {
	playerRepo := player.NewPostgresRepository(db)
	playerService := player.NewPlayerService(playerRepo)
	playerHandler := handler.NewPlayerHandler(playerService)

	router.HandleFunc("GET /player", playerHandler.GetAllPlayers)
	router.HandleFunc("GET /player/{id}", playerHandler.GetPlayerByID)
	router.HandleFunc("POST /player", playerHandler.CreatePlayer)
}
