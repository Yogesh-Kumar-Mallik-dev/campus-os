package http

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// BLOCK_HTTP_MIDDLEWARE_RECOVERY_001
// Purpose: Catches panics and converts them into structured 500 RFC 7807 responses.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("BLOCK_HTTP_MIDDLEWARE_RECOVERY_001: PANIC recovered: %v", rec)
				WriteProblem(w, r, http.StatusInternalServerError, "ERR_INTERNAL_PANIC", "Internal Server Error", "An unexpected server panic occurred", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// BLOCK_HTTP_MIDDLEWARE_TRACING_001
// Purpose: Attaches a unique request correlation ID to incoming requests and response headers.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = generateRequestID()
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)

		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Printf("BLOCK_HTTP_LOG_001: [%s] %s %s took %v", reqID, r.Method, r.URL.Path, time.Since(start))
	})
}

// BLOCK_HTTP_MIDDLEWARE_CORS_001
// Purpose: Configures standard CORS headers for Web & Desktop clients.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "req_fallback"
	}
	return "req_" + hex.EncodeToString(bytes)
}
