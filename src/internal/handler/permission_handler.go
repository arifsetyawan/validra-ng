package handler

import (
	"net/http"

	"github.com/arifsetyawan/validra/src/internal/dto"
	"github.com/arifsetyawan/validra/src/internal/service"
	"github.com/labstack/echo/v4"
)

// PermissionHandler handles HTTP requests related to permissions
type PermissionHandler struct {
	permissionService *service.PermissionService
}

// NewPermissionHandler creates a new PermissionHandler
func NewPermissionHandler(permissionService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// Register registers routes to the Echo instance
func (h *PermissionHandler) Register(e *echo.Echo) {
	// Update the path to include the /api prefix like other endpoints
	e.POST("/api/check-permission", h.CheckPermission)

	// Keep the original path for backward compatibility
	e.POST("/check-permission", h.CheckPermission)
}

// CheckPermission godoc
// @Summary Check permission
// @Description Checks if a user has permission to perform an action on a resource
// @Tags permissions
// @Accept json
// @Produce json
// @Param request body dto.PermissionCheckRequest true "Permission check request"
// @Success 200 {object} dto.PermissionCheckResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/check-permission [post]
func (h *PermissionHandler) CheckPermission(c echo.Context) error {
	// Bind request body to PermissionCheckRequest struct
	req := new(dto.PermissionCheckRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Check permission
	granted, context, err := h.permissionService.CheckPermission(c.Request().Context(), req.User, req.Action, req.Resource)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return response
	return c.JSON(http.StatusOK, dto.PermissionCheckResponse{
		Grant:   granted,
		Context: context,
	})
}
