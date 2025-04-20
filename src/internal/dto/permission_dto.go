package dto

// PermissionCheckRequest represents the request structure for checking permissions
type PermissionCheckRequest struct {
	User     string `json:"user" validate:"required"`
	Action   string `json:"action" validate:"required"`
	Resource string `json:"resource" validate:"required"`
}

// PermissionCheckResponse represents the response structure for permission check results
type PermissionCheckResponse struct {
	Grant   bool                   `json:"grant"`
	Context map[string]interface{} `json:"context"`
}
