package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

func TestUpdateWorkflowWithRetryRefreshesVersionAfterConflict(t *testing.T) {
	const workflowID = "11111111-1111-1111-1111-111111111111"
	putCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			version := 1
			if putCount > 0 {
				version = 2
			}
			writeWorkflow(t, writer, workflowID, version, false)
		case http.MethodPut:
			var update apiclient.UpdateWorkflowRequest
			if err := json.NewDecoder(request.Body).Decode(&update); err != nil {
				t.Fatalf("decode update request: %v", err)
			}
			wantVersion := putCount + 1
			if update.ExpectedVersion == nil || *update.ExpectedVersion != wantVersion {
				t.Fatalf("expected_version = %v, want %d", update.ExpectedVersion, wantVersion)
			}
			putCount++
			if putCount == 1 {
				writer.WriteHeader(http.StatusConflict)
				_, _ = writer.Write([]byte(`{"message":"stale workflow version"}`))
				return
			}
			writeWorkflow(t, writer, workflowID, 3, false)
		default:
			t.Fatalf("unexpected method %s", request.Method)
		}
	}))
	defer server.Close()

	client := apiclient.NewClient(server.URL, "token")
	updated, err := updateWorkflowWithRetry(client, workflowID, func(workflow *apiclient.Workflow) (apiclient.UpdateWorkflowRequest, error) {
		request := workflowUpdateRequest(workflow)
		return request, nil
	})
	if err != nil {
		t.Fatalf("updateWorkflowWithRetry() error = %v", err)
	}
	if updated.Version != 3 {
		t.Fatalf("updated workflow version = %d, want 3", updated.Version)
	}
	if putCount != 2 {
		t.Fatalf("PUT count = %d, want 2", putCount)
	}
}

func TestValidateExecutionIdentity(t *testing.T) {
	tests := []struct {
		name    string
		model   WorkflowExecutionIdentityResourceModel
		wantErr bool
	}{
		{
			name:    "service principal requires ID",
			model:   WorkflowExecutionIdentityResourceModel{Type: types.StringValue("service_principal")},
			wantErr: true,
		},
		{
			name: "service principal accepts ID",
			model: WorkflowExecutionIdentityResourceModel{
				Type:               types.StringValue("service_principal"),
				ServicePrincipalID: types.StringValue("22222222-2222-2222-2222-222222222222"),
			},
		},
		{
			name: "initiator does not accept ID",
			model: WorkflowExecutionIdentityResourceModel{
				Type:               types.StringValue("initiator"),
				ServicePrincipalID: types.StringValue("22222222-2222-2222-2222-222222222222"),
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateExecutionIdentity(test.model)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateExecutionIdentity() error = %v, wantErr %t", err, test.wantErr)
			}
		})
	}
}

func TestOptionalRawMessageStatePreservesJSONNull(t *testing.T) {
	got := optionalRawMessageState(json.RawMessage("null"))
	if got != types.StringValue("null") {
		t.Fatalf("optionalRawMessageState(null) = %#v, want string null", got)
	}
}

func writeWorkflow(t *testing.T, writer http.ResponseWriter, id string, version int, published bool) {
	t.Helper()
	workflow := apiclient.Workflow{
		ID:         id,
		Name:       "workflow",
		Version:    version,
		Published:  published,
		Conditions: json.RawMessage(`[]`),
		Steps:      json.RawMessage(`[]`),
		Graph:      json.RawMessage(`{}`),
	}
	if err := json.NewEncoder(writer).Encode(workflow); err != nil {
		t.Fatalf("encode workflow: %v", err)
	}
}
