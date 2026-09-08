package provider

import (
	"context"
	"os"
	"testing"

	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"moneat": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("MONEAT_AUTH_TOKEN"); v == "" {
		t.Fatal("MONEAT_AUTH_TOKEN must be set for acceptance tests")
	}
}

func TestProviderSchemaIncludesResponseAPIKey(t *testing.T) {
	var response frameworkprovider.SchemaResponse
	(&MoneatProvider{}).Schema(context.Background(), frameworkprovider.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("provider schema diagnostics: %v", response.Diagnostics)
	}
	if _, ok := response.Schema.Attributes["response_api_key"]; !ok {
		t.Fatal("provider schema does not include response_api_key")
	}
}
