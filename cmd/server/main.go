package main

import (
	"log"
	"net/http"

	router "github.com/icestormerrr/notes-api/internal/delivery/http"
	"github.com/icestormerrr/notes-api/internal/delivery/http/handlers"
	"github.com/icestormerrr/notes-api/internal/repos"
	"github.com/icestormerrr/notes-api/internal/services"
	"github.com/icestormerrr/notes-api/internal/utils/config"
)

func main() {
	cfg := config.Load()

	notesInmemoryRepo := repos.NewNoteRepoMem()
	notesService := services.NewNotesService(notesInmemoryRepo)
	notesHandler := handlers.NewNoteHandler(notesService)
	mux := router.Build(notesHandler)

	log.Println("listening on", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
