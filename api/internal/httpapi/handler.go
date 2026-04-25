package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"trimora/internal/links"
	"trimora/internal/validate"
)

type Handler struct {
	service *links.Service
	baseURL string
	log     *slog.Logger
}

func NewHandler(service *links.Service, baseURL string, log *slog.Logger) *Handler {
	return &Handler{service: service, baseURL: strings.TrimRight(baseURL, "/"), log: log}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	link, err := h.service.Create(r.Context(), req.URL, req.Alias)
	if err != nil {
		h.respondCreateError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateResponse{
		Code:     link.Code,
		ShortURL: h.baseURL + "/" + link.Code,
		URL:      link.TargetURL,
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	target, err := h.service.Resolve(r.Context(), code)
	if errors.Is(err, links.ErrNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		h.log.Error("resolve link", "code", code, "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (h *Handler) respondCreateError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, links.ErrAliasUnavailable):
		writeError(w, http.StatusConflict, err.Error())
	case isValidationError(err):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		h.log.Error("create link", "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func isValidationError(err error) bool {
	switch {
	case errors.Is(err, validate.ErrURLEmpty),
		errors.Is(err, validate.ErrURLTooLong),
		errors.Is(err, validate.ErrURLInvalid),
		errors.Is(err, validate.ErrURLScheme),
		errors.Is(err, validate.ErrAliasLength),
		errors.Is(err, validate.ErrAliasFormat),
		errors.Is(err, validate.ErrAliasReserved):
		return true
	}
	return false
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}
