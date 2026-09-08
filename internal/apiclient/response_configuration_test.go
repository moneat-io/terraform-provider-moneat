package apiclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponseConfigurationClientRequests(t *testing.T) {
	const configuration = `{"schemaVersion":1,"source":"TERRAFORM","resources":[]}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization header = %q, want bearer token", got)
		}
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/response/configuration/state":
			_, _ = io.WriteString(
				w,
				`{"id":"revision-id","revision":3,"checksum":"checksum-3","source":"TERRAFORM","configuration":`+
					configuration+`}`,
			)
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/response/configuration/plan":
			var request struct {
				Configuration    json.RawMessage `json:"configuration"`
				ExpectedRevision *int            `json:"expectedRevision"`
				AllowDestructive bool            `json:"allowDestructive"`
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("reading plan request: %v", err)
			}
			if err := json.Unmarshal(body, &request); err != nil {
				t.Errorf("decoding plan request: %v", err)
			}
			if string(request.Configuration) != configuration {
				t.Errorf("configuration = %s, want %s", request.Configuration, configuration)
			}
			if request.ExpectedRevision == nil || *request.ExpectedRevision != 3 {
				t.Errorf("expectedRevision = %v, want 3", request.ExpectedRevision)
			}
			if request.AllowDestructive {
				t.Error("allowDestructive = true, want false")
			}
			_, _ = io.WriteString(
				w,
				`{"baseRevision":3,"proposedRevision":4,"checksum":"checksum-4","changes":[],"issues":[],`+
					`"canApply":true,"configuration":`+configuration+`}`,
			)
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/response/configuration/apply":
			if got := r.Header.Get("Idempotency-Key"); got != "apply-key" {
				t.Errorf("Idempotency-Key = %q, want apply-key", got)
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("reading apply request: %v", err)
			}
			if !strings.Contains(string(body), `"idempotencyKey":"apply-key"`) {
				t.Errorf("apply request does not contain idempotency key: %s", body)
			}
			_, _ = io.WriteString(
				w,
				`{"id":"revision-id-4","revision":4,"checksum":"checksum-4","source":"TERRAFORM","configuration":`+
					configuration+`}`,
			)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	state, err := client.GetResponseConfigurationState()
	if err != nil {
		t.Fatalf("GetResponseConfigurationState() error = %v", err)
	}
	if state.Revision != 3 || state.ID != "revision-id" {
		t.Fatalf("state = %+v, want revision 3 and revision-id", state)
	}

	expectedRevision := 3
	plan, err := client.PlanResponseConfiguration(json.RawMessage(configuration), &expectedRevision, false)
	if err != nil {
		t.Fatalf("PlanResponseConfiguration() error = %v", err)
	}
	if !plan.CanApply || plan.ProposedRevision != 4 {
		t.Fatalf("plan = %+v, want applicable revision 4", plan)
	}

	applied, err := client.ApplyResponseConfiguration(plan, "apply-key", false)
	if err != nil {
		t.Fatalf("ApplyResponseConfiguration() error = %v", err)
	}
	if applied.Revision != 4 || applied.Checksum != "checksum-4" {
		t.Fatalf("applied = %+v, want revision 4 and checksum-4", applied)
	}
}

func TestApplyResponseConfigurationRejectsNilPlan(t *testing.T) {
	client := NewClient("http://example.invalid", "test-token")
	if _, err := client.ApplyResponseConfiguration(nil, "apply-key", false); err == nil {
		t.Fatal("ApplyResponseConfiguration(nil) returned nil error")
	}
}

func TestResponseConfigurationClientUsesDedicatedToken(t *testing.T) {
	client := NewClient("https://api.example.test", "primary-token")
	client.ResponseToken = "response-token"

	responseClient := client.ResponseConfigurationClient()
	if responseClient.Token != "response-token" {
		t.Fatalf("response client token = %q, want response-token", responseClient.Token)
	}
	if responseClient.BaseURL != client.BaseURL || responseClient.HTTPClient != client.HTTPClient {
		t.Fatal("response client did not preserve base URL and HTTP transport")
	}
}
