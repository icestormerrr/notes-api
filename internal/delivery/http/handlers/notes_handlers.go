package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/icestormerrr/notes-api/internal/core"
	http_utils "github.com/icestormerrr/notes-api/internal/utils/http"
)

type NoteHandler struct {
	service core.NotesService
}

func NewNoteHandler(service core.NotesService) *NoteHandler {
	return &NoteHandler{service: service}
}

type NoteCreateRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteUpdateRequest struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.GetAll()
	if err != nil {
		http_utils.WriteError(w, http.StatusInternalServerError, "failed to get notes", err.Error())
		return
	}

	http_utils.WriteJSON(w, notes)
}

func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid note id", err.Error())
		return
	}

	note, err := h.service.GetById(id)
	if err != nil {
		http_utils.WriteError(w, http.StatusNotFound, "note not found", nil)
		return
	}

	http_utils.WriteJSON(w, note)
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req NoteCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	payload := core.NoteCreatePayload{
		Title:   req.Title,
		Content: req.Content,
	}

	id, err := h.service.CreateNote(payload)
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	http_utils.WriteJSON(w, map[string]int64{"id": id})
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid note id", err.Error())
		return
	}

	var req NoteUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	payload := core.NoteUpdatePayload{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.service.UpdateNote(id, payload); err != nil {
		if err.Error() == "note not found" {
			http_utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		http_utils.WriteError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(chi.URLParam(r, "id"))
	if err != nil {
		http_utils.WriteError(w, http.StatusBadRequest, "invalid note id", err.Error())
		return
	}

	if err := h.service.Delete(id); err != nil {
		http_utils.WriteError(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(idStr string) (int64, error) {
	return strconv.ParseInt(idStr, 10, 64)
}
