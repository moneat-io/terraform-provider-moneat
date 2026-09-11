package apiclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWorkflowClientUsesUUIDPathsAndExpectedVersion(t *testing.T) {
	const workflowID = "7f9f3d3c-72b8-4f51-9a7d-2be0b4f6f8d0"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", request.Method)
		}
		wantPath := "/v1/workflows/" + workflowID
		if request.URL.Path != wantPath {
			t.Fatalf("path = %s, want %s", request.URL.Path, wantPath)
		}
		var payload UpdateWorkflowRequest
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload.ExpectedVersion == nil || *payload.ExpectedVersion != 4 {
			t.Fatalf("expected version = %v, want 4", payload.ExpectedVersion)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"id":"` + workflowID + `","name":"demo","trigger_name":"event","enabled":true,"version":5,"published":false,"conditions":[],"steps":[],"graph":{"nodes":[],"edges":[]},"once_for_template":[]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	version := 4
	workflow, err := client.UpdateWorkflow(workflowID, UpdateWorkflowRequest{ExpectedVersion: &version})
	if err != nil {
		t.Fatalf("UpdateWorkflow: %v", err)
	}
	if workflow.ID != workflowID {
		t.Fatalf("workflow ID = %q, want %q", workflow.ID, workflowID)
	}
}

func TestWorkflowClientUsesUUIDPublishQuery(t *testing.T) {
	const workflowID = "7f9f3d3c-72b8-4f51-9a7d-2be0b4f6f8d0"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v1/workflows/"+workflowID+"/publish" || request.URL.Query().Get("expected_version") != "5" {
			t.Fatalf("publish URL = %s, want UUID path with expected_version=5", request.URL.String())
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"id":"` + workflowID + `","version":6,"published":true}`))
	}))
	defer server.Close()

	version := 5
	workflow, err := NewClient(server.URL, "token").PublishWorkflow(workflowID, &version)
	if err != nil {
		t.Fatalf("PublishWorkflow: %v", err)
	}
	if !workflow.Published {
		t.Fatal("published = false, want true")
	}
}
