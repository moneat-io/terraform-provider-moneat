package provider

import (
	"encoding/json"
	"testing"
)

func TestResponseConfigurationStateStringOmitsGeneratedAt(t *testing.T) {
	got := responseConfigurationStateString(json.RawMessage(
		`{"generatedAt":"2026-09-08T12:00:00Z","resources":[],"schemaVersion":1}`,
	))
	want := `{"resources":[],"schemaVersion":1}`
	if got != want {
		t.Fatalf("responseConfigurationStateString() = %s, want %s", got, want)
	}
}
