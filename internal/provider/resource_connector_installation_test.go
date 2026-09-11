package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/moneat-io/terraform-provider-moneat/internal/apiclient"
)

func TestParseExternalAccountRejectsUnknownFields(t *testing.T) {
	_, err := parseExternalAccount(types.StringValue(`{"customerID":"acct-123"}`))
	if err == nil {
		t.Fatal("parseExternalAccount() accepted an unknown selector field")
	}
}

func TestImportedExternalAccountReconstructsProjectSelector(t *testing.T) {
	tests := []struct {
		name       string
		providerID string
		want       string
	}{
		{name: "project based connector", providerID: "pagerduty", want: `{"projectId":"acct-123"}`},
		{name: "google ads", providerID: "google_ads", want: `{"customerId":"acct-123"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := importedExternalAccount(&apiclient.ConnectorInstallation{
				ProviderID:        test.providerID,
				ExternalProjectID: "acct-123",
			})
			if got != types.StringValue(test.want) {
				t.Fatalf("imported external account = %#v, want %#v", got, types.StringValue(test.want))
			}
		})
	}
}

func TestImportedExternalAccountIsNullWithoutExternalProject(t *testing.T) {
	got := importedExternalAccount(&apiclient.ConnectorInstallation{})
	if !got.IsNull() {
		t.Fatalf("imported external account = %#v, want null", got)
	}
}
