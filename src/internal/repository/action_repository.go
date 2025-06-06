package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arifsetyawan/validra/src/internal/domain"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/google/uuid"
)

// GormActionRepository implements domain.ActionRepository using GORM with PostgreSQL
type ActionRepository struct {
	db *database.PostgresDB
}

// NewActionRepository creates a new GORM repository for actions
func NewActionRepository(db *database.PostgresDB) domain.ActionRepository {
	return &ActionRepository{
		db: db,
	}
}

// Action is the GORM model for actions
type Action struct {
	ID          string `gorm:"primaryKey"`
	ResourceID  string `gorm:"not null;index"`
	Name        string `gorm:"not null"`
	Description string
	Attributes  []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

// toDomain converts a GORM model to a domain model
func (a *Action) toDomain() *domain.Action {
	return &domain.Action{
		ID:          a.ID,
		ResourceID:  a.ResourceID,
		Name:        a.Name,
		Description: a.Description,
		Attributes:  a.Attributes,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// fromDomain converts a domain model to a GORM model
func actionFromDomain(a *domain.Action) *Action {
	return &Action{
		ID:          a.ID,
		ResourceID:  a.ResourceID,
		Name:        a.Name,
		Description: a.Description,
		Attributes:  a.Attributes,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// Create inserts a new action into the database
func (r *ActionRepository) Create(ctx context.Context, action *domain.Action) error {
	// Generate a new UUID if not provided
	if action.ID == "" {
		action.ID = uuid.New().String()
	}

	now := time.Now()
	action.CreatedAt = now
	action.UpdatedAt = now

	gormAction := actionFromDomain(action)
	result := r.db.DB.WithContext(ctx).Create(gormAction)
	if result.Error != nil {
		return fmt.Errorf("failed to create action: %w", result.Error)
	}

	return nil
}

// GetByID retrieves an action by ID
func (r *ActionRepository) GetByID(ctx context.Context, id string) (*domain.Action, error) {
	var action Action
	result := r.db.DB.WithContext(ctx).First(&action, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get action: %w", result.Error)
	}

	return action.toDomain(), nil
}

// GetByResourceID retrieves actions by resource ID
func (r *ActionRepository) GetByResourceID(ctx context.Context, resourceID string) ([]*domain.Action, error) {
	var actions []Action
	result := r.db.DB.WithContext(ctx).Where("resource_id = ?", resourceID).Find(&actions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get actions: %w", result.Error)
	}

	domainActions := make([]*domain.Action, len(actions))
	for i, action := range actions {
		domainActions[i] = action.toDomain()
	}

	return domainActions, nil
}

// List retrieves a paginated list of actions
func (r *ActionRepository) List(ctx context.Context, limit, offset int) ([]*domain.Action, error) {
	var actions []Action
	result := r.db.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&actions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list actions: %w", result.Error)
	}

	domainActions := make([]*domain.Action, len(actions))
	for i, action := range actions {
		domainActions[i] = action.toDomain()
	}

	return domainActions, nil
}

// Update updates an action in the database
func (r *ActionRepository) Update(ctx context.Context, action *domain.Action) error {
	action.UpdatedAt = time.Now()

	gormAction := actionFromDomain(action)
	result := r.db.DB.WithContext(ctx).Save(gormAction)
	if result.Error != nil {
		return fmt.Errorf("failed to update action: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("action not found")
	}

	return nil
}

// Delete removes an action from the database
func (r *ActionRepository) Delete(ctx context.Context, id string) (*domain.Action, error) {
	// First retrieve the action to return it after deletion
	var action Action
	getResult := r.db.DB.WithContext(ctx).First(&action, "id = ?", id)
	if getResult.Error != nil {
		return nil, fmt.Errorf("failed to get resource: %w", getResult.Error)
	}

	// Perform soft delete
	now := time.Now()
	result := r.db.DB.WithContext(ctx).Model(&Action{}).Where("id = ?", id).Update("deleted_at", now)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to soft delete action: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("action not found")
	}

	// Update the retrieved action with deletion time
	action.DeletedAt = &now

	return action.toDomain(), nil
}
