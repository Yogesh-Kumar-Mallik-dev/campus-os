package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BLOCK_RESPONSE_TEST_001
// Purpose: Verifies RFC 7807 problem details and collection empty semantics.
func TestWriteProblem(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/123", nil)
	w := httptest.NewRecorder()

	WriteProblem(w, req, http.StatusNotFound, "ERR_STUDENT_NOT_FOUND", "Student Not Found", "BLOCK_STUDENT_001: Student ID 123 does not exist", nil)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/problem+json; charset=utf-8" {
		t.Errorf("expected Content-Type application/problem+json, got %s", ct)
	}

	var problem ProblemDetails
	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		t.Fatalf("failed to decode problem json: %v", err)
	}
	if problem.Code != "ERR_STUDENT_NOT_FOUND" {
		t.Errorf("expected code ERR_STUDENT_NOT_FOUND, got %s", problem.Code)
	}
}

func TestWriteCollectionEmptySemantics(t *testing.T) {
	w := httptest.NewRecorder()
	var emptySlice []string

	WriteCollection(w, emptySlice)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK for empty collection, got %d", resp.StatusCode)
	}

	var res CollectionResponse[string]
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode collection response: %v", err)
	}
	if res.Data == nil || len(res.Data) != 0 {
		t.Errorf("expected non-nil empty slice in data, got %v", res.Data)
	}
	if res.TotalItems != 0 {
		t.Errorf("expected totalItems 0, got %d", res.TotalItems)
	}
}
