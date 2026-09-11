package apiclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConnectorClientUsesCurrentRoutesAndDoesNotDecodeSecret(t *testing.T) {
	const installationID = "2d5aab2c-fd2e-4d0e-8c6b-9aebf8ca18d1"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v1/connectors/installations/"+installationID+"/resources" {
			t.Fatalf("path = %s, want current connector resource route", request.URL.Path)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(`{"resources":[{"id":"resource-1","installationId":"` + installationID + `","externalResourceType":"channel","externalResourceId":"C123","displayName":"Alerts","providerMetadata":{}}]}`))
	}))
	defer server.Close()

	resources, err := NewClient(server.URL, "token").ListConnectorResources(installationID)
	if err != nil {
		t.Fatalf("ListConnectorResources: %v", err)
	}
	if len(resources) != 1 || resources[0].ExternalResourceID != "C123" {
		t.Fatalf("resources = %#v", resources)
	}
}
