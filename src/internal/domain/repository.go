package domain

import "context"

// ResourceRepository defines the methods for Resource data access
type ResourceRepository interface {
	Create(ctx context.Context, resource *Resource) error
	GetByID(ctx context.Context, id string) (*Resource, error)
	List(ctx context.Context, limit, offset int) ([]*Resource, error)
	Update(ctx context.Context, resource *Resource) error
	Delete(ctx context.Context, id string) error
}

// ActionRepository defines the methods for Action data access
type ActionRepository interface {
	Create(ctx context.Context, action *Action) error
	GetByID(ctx context.Context, id string) (*Action, error)
	GetByResourceID(ctx context.Context, resourceID string) ([]*Action, error)
	List(ctx context.Context, limit, offset int) ([]*Action, error)
	Update(ctx context.Context, action *Action) error
	Delete(ctx context.Context, id string) error
}

// RoleRepository defines the methods for Role data access
type RoleRepository interface {
	Create(ctx context.Context, role *Role) error
	GetByID(ctx context.Context, id string) (*Role, error)
	List(ctx context.Context, limit, offset int) ([]*Role, error)
	Update(ctx context.Context, role *Role) error
	Delete(ctx context.Context, id string) error
}

// UserRepository defines the methods for User data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, limit, offset int) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// UserSetRepository defines the methods for UserSet data access
type UserSetRepository interface {
	Create(ctx context.Context, userSet *UserSet) error
	GetByID(ctx context.Context, id string) (*UserSet, error)
	List(ctx context.Context, limit, offset int) ([]*UserSet, error)
	Update(ctx context.Context, userSet *UserSet) error
	Delete(ctx context.Context, id string) error
}

// ResourceSetRepository defines the methods for ResourceSet data access
type ResourceSetRepository interface {
	Create(ctx context.Context, resourceSet *ResourceSet) error
	GetByID(ctx context.Context, id string) (*ResourceSet, error)
	List(ctx context.Context, limit, offset int) ([]*ResourceSet, error)
	Update(ctx context.Context, resourceSet *ResourceSet) error
	Delete(ctx context.Context, id string) error
}

// PermissionRepository defines the methods for Permission data access
type PermissionRepository interface {
	Create(ctx context.Context, permission *Permission) error
	GetByID(ctx context.Context, id string) (*Permission, error)
	List(ctx context.Context, limit, offset int) ([]*Permission, error)
	Update(ctx context.Context, permission *Permission) error
	Delete(ctx context.Context, id string) error
	CheckPermission(ctx context.Context, userID, resourceID string) (bool, error)
}
