package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func normalizeJSONString(value string) (string, error) {
	var compacted bytes.Buffer
	if err := json.Compact(&compacted, []byte(value)); err != nil {
		return "", fmt.Errorf("invalid JSON: %w", err)
	}
	return compacted.String(), nil
}

func rawMessageFromJSONString(value string) (json.RawMessage, string, error) {
	normalized, err := normalizeJSONString(value)
	if err != nil {
		return nil, "", err
	}
	return json.RawMessage(normalized), normalized, nil
}

func rawMessageString(value json.RawMessage) string {
	if len(value) == 0 {
		return ""
	}
	normalized, err := normalizeJSONString(string(value))
	if err != nil {
		return string(value)
	}
	return normalized
}

func responseConfigurationStateString(value json.RawMessage) string {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(value, &object); err != nil {
		return rawMessageString(value)
	}
	delete(object, "generatedAt")
	normalized, err := json.Marshal(object)
	if err != nil {
		return rawMessageString(value)
	}
	return string(normalized)
}

func parseTerraformID(value types.String, label string) (int, error) {
	id, err := strconv.Atoi(value.ValueString())
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", label, err)
	}
	return id, nil
}

func terraformID(value int) types.String {
	return types.StringValue(strconv.Itoa(value))
}
