package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Router builds the HTTP router for the API.
func Router(h *Handler, allowedOrigins []string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(cors(allowedOrigins))

	r.Get("/healthz", h.Live)
	r.Get("/livez", h.Live)
	r.Get("/readyz", h.Ready)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", h.Create)
	})

	r.Get("/{code:[A-Za-z0-9_-]+}", h.Redirect)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		if wantsHTML(req) {
			renderStatusPage(w, http.StatusNotFound, "Page not found",
				"That page doesn’t exist on Trimora.")
			return
		}
		writeError(w, http.StatusNotFound, "not found")
	})

	return r
}

func cors(origins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		allowed[o] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowed[origin]; ok || len(origins) == 0 {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				}
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
