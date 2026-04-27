package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"trimora/internal/links"
	"trimora/internal/validate"
)

type Handler struct {
	service *links.Service
	db      *sql.DB
	baseURL string
	log     *slog.Logger
}

func NewHandler(service *links.Service, db *sql.DB, baseURL string, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		db:      db,
		baseURL: strings.TrimRight(baseURL, "/"),
		log:     log,
	}
}

// Live reports that the process is up. It does not depend on any backing service.
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready reports that the process can serve traffic, including the database.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.db.PingContext(ctx); err != nil {
		h.log.Warn("readiness check failed", "err", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "unavailable",
			"error":  "database unreachable",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	expiresIn, err := validate.ExpiresIn(req.ExpiresIn)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	link, err := h.service.Create(r.Context(), req.URL, req.Alias, expiresIn)
	if err != nil {
		h.respondCreateError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreateResponse{
		Code:      link.Code,
		ShortURL:  h.baseURL + "/" + link.Code,
		URL:       link.TargetURL,
		ExpiresAt: link.ExpiresAt,
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	target, err := h.service.Resolve(r.Context(), code)
	switch {
	case errors.Is(err, links.ErrNotFound):
		h.notFound(w, r)
		return
	case errors.Is(err, links.ErrExpired):
		h.gone(w, r)
		return
	case err != nil:
		h.log.Error("resolve link", "code", code, "err", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	if wantsHTML(r) {
		renderStatusPage(w, http.StatusNotFound, "Link not found",
			"That short link doesn’t exist or was never created.")
		return
	}
	writeError(w, http.StatusNotFound, "link not found")
}

func (h *Handler) gone(w http.ResponseWriter, r *http.Request) {
	if wantsHTML(r) {
		renderStatusPage(w, http.StatusGone, "Link expired",
			"This short link has reached its expiry and is no longer active.")
		return
	}
	writeError(w, http.StatusGone, "link expired")
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
		errors.Is(err, validate.ErrAliasReserved),
		errors.Is(err, validate.ErrExpiryInvalid):
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

// wantsHTML returns true when the client prefers an HTML response, e.g. a
// browser following a short link.
func wantsHTML(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if accept == "" {
		return false
	}
	if strings.Contains(accept, "application/json") {
		return false
	}
	return strings.Contains(accept, "text/html") || strings.Contains(accept, "*/*")
}

type statusPageData struct {
	Status  int
	Title   string
	Message string
}

var statusPageTpl = template.Must(template.New("status").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<title>{{.Title}} · trimora</title>
<style>
  :root {
    color-scheme: light;
    --bg: #f6f1ea;
    --card: #ffffff;
    --ink: #1f1b16;
    --muted: #6b6258;
    --accent: #c0623f;
  }
  * { box-sizing: border-box; }
  html, body { margin: 0; padding: 0; height: 100%; }
  body {
    background: var(--bg);
    color: var(--ink);
    font-family: ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, sans-serif;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px;
  }
  main {
    width: 100%;
    max-width: 460px;
    background: var(--card);
    border-radius: 16px;
    padding: 32px 28px;
    box-shadow: 0 12px 30px rgba(40, 28, 18, 0.08);
    text-align: center;
  }
  .mark {
    font-size: 14px;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: var(--muted);
    margin-bottom: 18px;
  }
  h1 {
    font-size: 26px;
    margin: 0 0 12px;
    line-height: 1.2;
  }
  p {
    margin: 0 0 24px;
    color: var(--muted);
    line-height: 1.5;
  }
  a.button {
    display: inline-block;
    padding: 12px 22px;
    border-radius: 999px;
    background: var(--accent);
    color: #fff;
    text-decoration: none;
    font-weight: 600;
  }
  .code {
    display: block;
    margin-top: 20px;
    font-size: 12px;
    color: var(--muted);
    letter-spacing: 0.1em;
  }
</style>
</head>
<body>
  <main>
    <div class="mark">trimora</div>
    <h1>{{.Title}}</h1>
    <p>{{.Message}}</p>
    <a class="button" href="/">Shorten a link</a>
    <span class="code">Status {{.Status}}</span>
  </main>
</body>
</html>`))

func renderStatusPage(w http.ResponseWriter, status int, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = statusPageTpl.Execute(w, statusPageData{
		Status:  status,
		Title:   title,
		Message: message,
	})
}
