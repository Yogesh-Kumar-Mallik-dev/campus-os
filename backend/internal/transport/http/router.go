package http

import (
	"net/http"

	campusv1 "github.com/Yogesh-Kumar-Mallik-dev/campus-os/backend/pkg/proto/campus/v1"
)

// RouterConfig represents parameters needed to construct the HTTP router.
type RouterConfig struct {
	Version    string
	AuthClient campusv1.AuthServiceClient
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

	// Auth & Onboarding endpoints
	authHandler := NewAuthHTTPHandler(cfg.AuthClient)
	mux.HandleFunc("POST /api/v1/auth/claim/validate", authHandler.HandleValidateClaim)
	mux.HandleFunc("POST /api/v1/auth/claim/verify-sim", authHandler.HandleVerifySIM)
	mux.HandleFunc("POST /api/v1/auth/claim/complete", authHandler.HandleCompleteClaim)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.HandleLogin)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.HandleRefresh)

	// Wrap mux with standard middleware stack
	handler := RecoveryMiddleware(mux)
	handler = CORSMiddleware(handler)
	handler = RequestIDMiddleware(handler)

	return handler
}
