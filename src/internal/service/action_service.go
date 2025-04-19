package service

import (
	"context"
	"fmt"

	"github.com/arifsetyawan/validra/src/internal/domain"
)

// ActionService handles business logic for actions
type ActionService struct {
	actionRepo   domain.ActionRepository
	resourceRepo domain.ResourceRepository
}

// NewActionService creates a new ActionService
func NewActionService(actionRepo domain.ActionRepository, resourceRepo domain.ResourceRepository) *ActionService {
	return &ActionService{
		actionRepo:   actionRepo,
		resourceRepo: resourceRepo,
	}
}

// CreateAction creates a new action
func (s *ActionService) CreateAction(ctx context.Context, action *domain.Action) error {
	if action.Name == "" {
		return fmt.Errorf("action name is required")
	}

	if action.ResourceID == "" {
		return fmt.Errorf("resource ID is required")
	}

	// Verify that the referenced resource exists
	_, err := s.resourceRepo.GetByID(ctx, action.ResourceID)
	if err != nil {
		return fmt.Errorf("invalid resource ID: %w", err)
	}

	return s.actionRepo.Create(ctx, action)
}

// GetActionByID retrieves an action by ID
func (s *ActionService) GetActionByID(ctx context.Context, id string) (*domain.Action, error) {
	return s.actionRepo.GetByID(ctx, id)
}

// GetActionsByResourceID retrieves actions by resource ID
func (s *ActionService) GetActionsByResourceID(ctx context.Context, resourceID string) ([]*domain.Action, error) {
	return s.actionRepo.GetByResourceID(ctx, resourceID)
}

// ListActions retrieves a paginated list of actions
func (s *ActionService) ListActions(ctx context.Context, limit, offset int) ([]*domain.Action, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	return s.actionRepo.List(ctx, limit, offset)
}

// UpdateAction updates an existing action
func (s *ActionService) UpdateAction(ctx context.Context, action *domain.Action) error {
	if action.ID == "" {
		return fmt.Errorf("action ID is required")
	}
	if action.Name == "" {
		return fmt.Errorf("action name is required")
	}
	if action.ResourceID == "" {
		return fmt.Errorf("resource ID is required")
	}

	// Verify that the referenced resource exists
	_, err := s.resourceRepo.GetByID(ctx, action.ResourceID)
	if err != nil {
		return fmt.Errorf("invalid resource ID: %w", err)
	}

	return s.actionRepo.Update(ctx, action)
}

// DeleteAction deletes an action by ID
func (s *ActionService) DeleteAction(ctx context.Context, id string) (*domain.Action, error) {
	return s.actionRepo.Delete(ctx, id)
}

// ActionRepository returns the action repository
func (s *ActionService) ActionRepository() domain.ActionRepository {
	return s.actionRepo
}
