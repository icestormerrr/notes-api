package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/icestormerrr/notes-api/internal/delivery/http/handlers"
)

func Build(notesHandler *handlers.NoteHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Get("/", notesHandler.GetAllNotes)
		r.Post("/", notesHandler.CreateNote)
		r.Get("/{id}", notesHandler.GetNoteByID)
		r.Patch("/{id}", notesHandler.UpdateNote)
		r.Delete("/{id}", notesHandler.DeleteNote)
	})

	return r
}
