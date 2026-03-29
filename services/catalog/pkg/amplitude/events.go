package amplitude

// Event type constants for the CRM service.
const (
	// Workspace / Organisation
	EventWorkspaceCreated = "workspace_created"

	// Tags
	EventTagCreated = "tag_created"
)

// Property key constants — use these instead of raw strings to avoid typos.
const (
	PropUserID       = "user_id"
	PropWorkspaceID  = "workspace_id"
	PropCrmID        = "crm_id"
	PropTagID        = "tag_id"
	PropTagName      = "tag_name"
	PropBusinessName = "business_name"
)
