package service

import (
	"context"
	"fmt"

	"github.com/arifsetyawan/validra/src/internal/model"
	"github.com/arifsetyawan/validra/src/internal/repository"
)

// UserService handles business logic for users
type UserService struct {
	userRepo repository.UserRepositoryInterface
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	// Check if username is already taken
	existingUser, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return fmt.Errorf("username already exists")
	}

	return s.userRepo.Create(ctx, user)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) (*[]model.User, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	return s.userRepo.List(ctx, limit, offset)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		return fmt.Errorf("user ID is required")
	}
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	// Check if we are changing the username and if it conflicts with an existing one
	existingUser, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil && existingUser.ID != user.ID {
		return fmt.Errorf("username already exists")
	}

	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.Delete(ctx, id)
}

// UserRepository returns the user repository
func (s *UserService) UserRepository() repository.UserRepositoryInterface {
	return s.userRepo
}
