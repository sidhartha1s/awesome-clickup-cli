// Copyright 2026 sidhartha1s. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"encoding/json"
	"testing"
)

// Test isNilOrEmpty with various JSON payloads.
func TestIsNilOrEmpty_WithName(t *testing.T) {
	raw := json.RawMessage(`{"name": "Test Task", "id": "abc123"}`)
	if isNilOrEmpty(raw) {
		t.Error("expected false for object with name, got true")
	}
}

func TestIsNilOrEmpty_WithEmptyName(t *testing.T) {
	raw := json.RawMessage(`{"name": "", "id": ""}`)
	if !isNilOrEmpty(raw) {
		t.Error("expected true for object with empty name and id, got false")
	}
}

func TestIsNilOrEmpty_WithNumericID(t *testing.T) {
	raw := json.RawMessage(`{"id": 12345}`)
	if isNilOrEmpty(raw) {
		t.Error("expected false for object with numeric id, got true")
	}
}

func TestIsNilOrEmpty_WithNestedDocument(t *testing.T) {
	raw := json.RawMessage(`{"score": 0.9, "document": {"name": "Task Name"}}`)
	if isNilOrEmpty(raw) {
		t.Error("expected false for search result with document.name, got true")
	}
}

func TestIsNilOrEmpty_WithScoreOnly(t *testing.T) {
	raw := json.RawMessage(`{"score": 0.5}`)
	if isNilOrEmpty(raw) {
		t.Error("expected false for object with score field, got true")
	}
}

func TestIsNilOrEmpty_InvalidJSON(t *testing.T) {
	raw := json.RawMessage(`not valid json`)
	if !isNilOrEmpty(raw) {
		t.Error("expected true for invalid JSON, got false")
	}
}

// Test extractSearchResults with various response formats.
func TestExtractSearchResults_DirectArray(t *testing.T) {
	raw := json.RawMessage(`[{"id": "1"}, {"id": "2"}]`)
	results := extractSearchResults(raw)
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestExtractSearchResults_DataWrapper(t *testing.T) {
	raw := json.RawMessage(`{"data": [{"id": "1"}, {"id": "2"}]}`)
	results := extractSearchResults(raw)
	if len(results) != 2 {
		t.Errorf("expected 2 results from data wrapper, got %d", len(results))
	}
}

func TestExtractSearchResults_ResultsWrapper(t *testing.T) {
	raw := json.RawMessage(`{"results": [{"id": "1"}]}`)
	results := extractSearchResults(raw)
	if len(results) != 1 {
		t.Errorf("expected 1 result from results wrapper, got %d", len(results))
	}
}

func TestExtractSearchResults_ItemsWrapper(t *testing.T) {
	raw := json.RawMessage(`{"items": [{"id": "a"}, {"id": "b"}, {"id": "c"}]}`)
	results := extractSearchResults(raw)
	if len(results) != 3 {
		t.Errorf("expected 3 results from items wrapper, got %d", len(results))
	}
}

func TestExtractSearchResults_SingleObject(t *testing.T) {
	raw := json.RawMessage(`{"id": "single"}`)
	results := extractSearchResults(raw)
	if len(results) != 1 {
		t.Errorf("expected 1 result for single object, got %d", len(results))
	}
}
