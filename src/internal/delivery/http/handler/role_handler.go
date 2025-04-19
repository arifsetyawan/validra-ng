package handler

import (
	"net/http"
	"strconv"

	"github.com/arifsetyawan/validra/src/internal/delivery/http/dto"
	"github.com/arifsetyawan/validra/src/internal/service"
	"github.com/labstack/echo/v4"
)

// RoleHandler handles HTTP requests for roles
type RoleHandler struct {
	roleService *service.RoleService
}

// NewRoleHandler creates a new RoleHandler
func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// Register registers the routes to the given echo instance
func (h *RoleHandler) Register(e *echo.Echo) {
	roles := e.Group("/api/roles")
	roles.POST("", h.CreateRole)
	roles.GET("", h.ListRoles)
	roles.GET("/:id", h.GetRole)
	roles.PUT("/:id", h.UpdateRole)
	roles.DELETE("/:id", h.DeleteRole)
}

// CreateRole creates a new role
// @Summary Create a new role
// @Description Create a new role with the provided information
// @Tags roles
// @Accept json
// @Produce json
// @Param role body dto.CreateRoleRequest true "Role information"
// @Success 201 {object} dto.RoleResponse "Role created"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/roles [post]
func (h *RoleHandler) CreateRole(c echo.Context) error {
	var req dto.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	role := req.ToRoleDomain()
	if err := h.roleService.CreateRole(c.Request().Context(), role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToRoleResponse(role)
	return c.JSON(http.StatusCreated, response)
}

// GetRole retrieves a role by ID
// @Summary Get a role by ID
// @Description Retrieve a specific role by its unique identifier
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} dto.RoleResponse "Role found"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Role not found"
// @Router /api/roles/{id} [get]
func (h *RoleHandler) GetRole(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing role ID"})
	}

	role, err := h.roleService.GetRoleByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
	}

	response := dto.ToRoleResponse(role)
	return c.JSON(http.StatusOK, response)
}

// ListRoles retrieves a paginated list of roles
// @Summary List roles
// @Description Get a paginated list of all roles
// @Tags roles
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to return (default: 10)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Success 200 {object} dto.ListRolesResponse "List of roles"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/roles [get]
func (h *RoleHandler) ListRoles(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	roles, err := h.roleService.ListRoles(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert domain models to response DTOs
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, r := range roles {
		roleResponses[i] = dto.ToRoleResponse(r)
	}

	response := dto.ListRolesResponse{
		Roles: roleResponses,
		Total: len(roleResponses), // In a real app, we'd get the total count from the service
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateRole updates an existing role
// @Summary Update a role
// @Description Update an existing role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body dto.UpdateRoleRequest true "Updated role information"
// @Success 200 {object} dto.RoleResponse "Role updated"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing role ID"})
	}

	var req dto.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get existing role
	role, err := h.roleService.GetRoleByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
	}

	// Update the role with request data
	req.UpdateRoleDomain(role)

	if err := h.roleService.UpdateRole(c.Request().Context(), role); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToRoleResponse(role)
	return c.JSON(http.StatusOK, response)
}

// DeleteRole deletes a role by ID
// @Summary Delete a role
// @Description Delete a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} dto.RoleResponse "Role soft deleted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing role ID"})
	}

	deletedRole, err := h.roleService.DeleteRole(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "role not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Role not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToRoleResponse(deletedRole)
	return c.JSON(http.StatusOK, response)
}
