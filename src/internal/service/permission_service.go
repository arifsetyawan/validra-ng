package service

import (
	"context"

	"github.com/arifsetyawan/validra/src/internal/repository"
)

// PermissionService handles business logic for permission checking
type PermissionService struct {
	userRepo     repository.UserRepositoryInterface
	resourceRepo repository.ResourceRepositoryInterface
	roleRepo     repository.RoleRepositoryInterface
}

// NewPermissionService creates a new PermissionService
func NewPermissionService(
	userRepo repository.UserRepositoryInterface,
	resourceRepo repository.ResourceRepositoryInterface,
	roleRepo repository.RoleRepositoryInterface,
) *PermissionService {
	return &PermissionService{
		userRepo:     userRepo,
		resourceRepo: resourceRepo,
		roleRepo:     roleRepo,
	}
}

// CheckPermission checks if a user has permission to perform an action on a resource
func (s *PermissionService) CheckPermission(ctx context.Context, username, actionName, resourceName string) (bool, map[string]interface{}, error) {
	// Context contains additional information about the permission decision
	context := map[string]interface{}{
		"userName":     username,
		"actionName":   actionName,
		"resourceName": resourceName,
		"roles":        []string{},
	}

	// Find the user by username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		// For testing purposes, we'll create a dummy user with the requested name
		// In a production environment, you might want to return an error instead
		context["userId"] = "unknown"
		context["userExists"] = false
	} else {
		context["userId"] = user.ID
		context["userExists"] = true
	}

	// Find the resource by name (simplified - in production you would likely need a
	// more efficient method to find resources by name)
	resources, err := s.resourceRepo.List(ctx, 100, 0)
	if err == nil {
		// Try to find the resource
		var resourceFound bool
		for _, r := range *resources {
			if r.Name == resourceName {
				context["resourceId"] = r.ID
				resourceFound = true

				for _, a := range r.ResourceActions {
					if a.Name == actionName {
						context["actionId"] = a.ID
						context["actionExists"] = true
						break
					}
				}

				break
			}
		}

		if !resourceFound {
			context["resourceId"] = "unknown"
			context["resourceExists"] = false
		}
	}

	// If we didn't find an action previously, set appropriate context values
	if _, exists := context["actionId"]; !exists {
		context["actionId"] = "unknown"
		context["actionExists"] = false
	}

	// For this initial implementation, we'll return a simple grant
	// In a real implementation, this would check against the Permission repository
	// to determine if the user has the required permissions
	// Currently always allowing access for testing purposes
	return true, context, nil
}
