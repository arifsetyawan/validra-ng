package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arifsetyawan/validra/src/internal/model"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/google/uuid"
)

// ActionRepository is a GORM implementation of the ActionRepositoryInterface
type ActionRepository struct {
	db *database.PostgresDB
}

// ActionRepositoryInterface defines the methods for Action data access
type ActionRepositoryInterface interface {
	Create(ctx context.Context, action *model.Action) error
	GetByID(ctx context.Context, id string) (*model.Action, error)
	GetByResourceID(ctx context.Context, resourceID string) (*[]model.Action, error)
	List(ctx context.Context, limit, offset int) (*[]model.Action, error)
	Update(ctx context.Context, action *model.Action) error
	Delete(ctx context.Context, id string) (*model.Action, error)
}

// NewActionRepository creates a new GORM repository for actions
func NewActionRepository(db *database.PostgresDB) ActionRepositoryInterface {
	return &ActionRepository{
		db: db,
	}
}

// Create inserts a new action into the database
func (r *ActionRepository) Create(ctx context.Context, action *model.Action) error {
	// Generate a new UUID if not provided
	if action.ID == "" {
		action.ID = uuid.New().String()
	}

	now := time.Now()
	action.CreatedAt = now
	action.UpdatedAt = now

	result := r.db.DB.WithContext(ctx).Create(action)
	if result.Error != nil {
		return fmt.Errorf("failed to create action: %w", result.Error)
	}

	return nil
}

// GetByID retrieves an action by ID
func (r *ActionRepository) GetByID(ctx context.Context, id string) (*model.Action, error) {
	var action model.Action
	result := r.db.DB.WithContext(ctx).First(&action, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get action: %w", result.Error)
	}

	return &action, nil
}

// GetByResourceID retrieves actions by resource ID
func (r *ActionRepository) GetByResourceID(ctx context.Context, resourceID string) (*[]model.Action, error) {
	var actions []model.Action
	result := r.db.DB.WithContext(ctx).Where("resource_id = ?", resourceID).Find(&actions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get actions: %w", result.Error)
	}

	return &actions, nil
}

// List retrieves a paginated list of actions
func (r *ActionRepository) List(ctx context.Context, limit, offset int) (*[]model.Action, error) {
	var actions []model.Action
	result := r.db.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&actions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list actions: %w", result.Error)
	}

	return &actions, nil
}

// Update updates an action in the database
func (r *ActionRepository) Update(ctx context.Context, action *model.Action) error {
	action.UpdatedAt = time.Now()

	result := r.db.DB.WithContext(ctx).Save(action)
	if result.Error != nil {
		return fmt.Errorf("failed to update action: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("action not found")
	}

	return nil
}

// Delete removes an action from the database
func (r *ActionRepository) Delete(ctx context.Context, id string) (*model.Action, error) {
	// First retrieve the action to return it after deletion
	var action model.Action
	getResult := r.db.DB.WithContext(ctx).First(&action, "id = ?", id)
	if getResult.Error != nil {
		return nil, fmt.Errorf("failed to get resource: %w", getResult.Error)
	}

	// Perform soft delete
	now := time.Now()
	result := r.db.DB.WithContext(ctx).Model(&action).Where("id = ?", id).Update("deleted_at", now)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to soft delete action: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("action not found")
	}

	// Update the retrieved action with deletion time
	action.DeletedAt = &now

	return &action, nil
}
