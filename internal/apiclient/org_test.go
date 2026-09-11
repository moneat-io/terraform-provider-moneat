package apiclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrgMemberClientUsesUuidEnvelopeAndRoleRoute(t *testing.T) {
	const memberID = "2d5aab2c-fd2e-4d0e-8c6b-9aebf8ca18d1"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodGet && request.URL.Path == "/v1/org/members":
			_, _ = writer.Write([]byte(`{"members":[{"userId":"` + memberID + `","email":"member@example.com","role":"member"}],"pendingInvitations":[]}`))
		case request.Method == http.MethodPut && request.URL.Path == "/v1/org/members/"+memberID+"/role":
			var body UpdateOrgMemberRequest
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatalf("decode role request: %v", err)
			}
			if body.Role != "admin" {
				t.Fatalf("role = %q, want admin", body.Role)
			}
			_, _ = writer.Write([]byte(`{"success":true}`))
		case request.Method == http.MethodDelete && request.URL.Path == "/v1/org/members/"+memberID:
			_, _ = writer.Write([]byte(`{"success":true}`))
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "token")
	members, err := client.ListOrgMembers()
	if err != nil {
		t.Fatalf("ListOrgMembers: %v", err)
	}
	if len(members) != 1 || members[0].ID != memberID {
		t.Fatalf("members = %#v", members)
	}
	if err := client.UpdateOrgMember(memberID, UpdateOrgMemberRequest{Role: "admin"}); err != nil {
		t.Fatalf("UpdateOrgMember: %v", err)
	}
	if err := client.DeleteOrgMember(memberID); err != nil {
		t.Fatalf("DeleteOrgMember: %v", err)
	}
}
