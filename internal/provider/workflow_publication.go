package provider

import "github.com/moneat-io/terraform-provider-moneat/internal/apiclient"

// preserveWorkflowPublication keeps the caller's publication intent when a
// companion resource (schedule or execution identity) creates a new version.
// Moneat intentionally makes every configuration change a draft, so a
// secondary Terraform resource must complete the same publish transition as
// the primary workflow resource. The intent is passed explicitly so a retry
// does not infer it from an already-drafted response.
func preserveWorkflowPublication(
	client *apiclient.Client,
	previous *apiclient.Workflow,
	updated *apiclient.Workflow,
	publicationIntent bool,
) (*apiclient.Workflow, error) {
	if !publicationIntent {
		return updated, nil
	}
	version := updated.Version
	return client.PublishWorkflow(updated.ID, &version)
}

// workflowUpdateRequest copies the complete mutable workflow payload before a
// companion Terraform resource changes one field. This prevents a partial
// update from clearing definition, identity, or governance-related fields.
func workflowUpdateRequest(workflow *apiclient.Workflow) apiclient.UpdateWorkflowRequest {
	enabled := workflow.Enabled
	version := workflow.Version
	return apiclient.UpdateWorkflowRequest{
		Name:              workflow.Name,
		Enabled:           &enabled,
		Conditions:        workflow.Conditions,
		Steps:             workflow.Steps,
		Graph:             workflow.Graph,
		OnceForTemplate:   workflow.OnceForTemplate,
		ExpectedVersion:   &version,
		RunOnce:           workflow.RunOnce,
		InputSchema:       workflow.InputSchema,
		TriggerNames:      workflow.TriggerNames,
		Schedules:         workflow.Schedules,
		ExecutionIdentity: workflow.ExecutionIdentity,
	}
}
