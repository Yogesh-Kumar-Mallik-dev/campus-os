package http

import (
	"encoding/json"
	"net/http"
)

// ProblemDetails represents an RFC 7807 standardized error envelope.
type ProblemDetails struct {
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	Status        int            `json:"status"`
	Detail        string         `json:"detail"`
	Instance      string         `json:"instance"`
	Code          string         `json:"code"`
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// BLOCK_HTTP_RESPONSE_JSON_001
// Purpose: Writes a standard JSON response payload with Content-Type header.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// BLOCK_HTTP_RESPONSE_ERROR_001
// Purpose: Writes an RFC 7807 standardized problem details error response.
func WriteProblem(w http.ResponseWriter, r *http.Request, status int, code, title, detail string, invalidParams []InvalidParam) {
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(status)

	problem := ProblemDetails{
		Type:          "https://campus-os.internal/errors/" + code,
		Title:         title,
		Status:        status,
		Detail:        detail,
		Instance:      r.URL.Path,
		Code:          code,
		InvalidParams: invalidParams,
	}

	_ = json.NewEncoder(w).Encode(problem)
}

// CollectionResponse standardizes collection queries ensuring HTTP 200 with data array.
type CollectionResponse[T any] struct {
	Data       []T `json:"data"`
	TotalItems int `json:"totalItems"`
}

// BLOCK_HTTP_RESPONSE_COLLECTION_001
// Purpose: Renders collection results guaranteeing empty array semantics (never 404).
func WriteCollection[T any](w http.ResponseWriter, items []T) {
	if items == nil {
		items = make([]T, 0)
	}
	WriteJSON(w, http.StatusOK, CollectionResponse[T]{
		Data:       items,
		TotalItems: len(items),
	})
}
