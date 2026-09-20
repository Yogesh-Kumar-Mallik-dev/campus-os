package http

import (
	"net/http"
)

// RouterConfig represents parameters needed to construct the HTTP router.
type RouterConfig struct {
	Version string
}

// BLOCK_HTTP_ROUTER_INIT_001
// Purpose: Constructs the HTTP ServeMux routing tree with standard middleware wrapping.
func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	// Health & readiness endpoints
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]any{
			"status":  "healthy",
			"service": "campus-os-backend",
			"version": cfg.Version,
		})
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]any{
			"ready": true,
		})
	})

	// Base API info
	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]any{
			"api":       "Campus OS Institutional API",
			"version":   "v1",
			"standards": "RFC 7807, Scoped RBAC, Protobuf Persistence over UDS",
		})
	})

	// Wrap mux with standard middleware stack
	handler := RecoveryMiddleware(mux)
	handler = CORSMiddleware(handler)
	handler = RequestIDMiddleware(handler)

	return handler
}
