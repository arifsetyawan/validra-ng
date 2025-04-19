package handler

import (
	"net/http"
	"strconv"

	"github.com/arifsetyawan/validra/src/internal/delivery/http/dto"
	"github.com/arifsetyawan/validra/src/internal/service"
	"github.com/labstack/echo/v4"
)

// ActionHandler handles HTTP requests for actions
type ActionHandler struct {
	actionService *service.ActionService
}

// NewActionHandler creates a new ActionHandler
func NewActionHandler(actionService *service.ActionService) *ActionHandler {
	return &ActionHandler{
		actionService: actionService,
	}
}

// Register registers the routes to the given echo instance
func (h *ActionHandler) Register(e *echo.Echo) {
	actions := e.Group("/api/actions")
	actions.POST("", h.CreateAction)
	actions.GET("", h.ListActions)
	actions.GET("/:id", h.GetAction)
	actions.PUT("/:id", h.UpdateAction)
	actions.DELETE("/:id", h.DeleteAction)
	actions.GET("/resource/:resourceID", h.GetActionsByResourceID)
}

// CreateAction creates a new action
// @Summary Create a new action
// @Description Create a new action with the provided information
// @Tags actions
// @Accept json
// @Produce json
// @Param action body dto.CreateActionRequest true "Action information"
// @Success 201 {object} dto.ActionResponse "Action created"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/actions [post]
func (h *ActionHandler) CreateAction(c echo.Context) error {
	var req dto.CreateActionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	action := req.ToActionDomain()
	if err := h.actionService.CreateAction(c.Request().Context(), action); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToActionResponse(action)
	return c.JSON(http.StatusCreated, response)
}

// GetAction retrieves an action by ID
// @Summary Get an action by ID
// @Description Retrieve a specific action by its unique identifier
// @Tags actions
// @Accept json
// @Produce json
// @Param id path string true "Action ID"
// @Success 200 {object} dto.ActionResponse "Action found"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Action not found"
// @Router /api/actions/{id} [get]
func (h *ActionHandler) GetAction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing action ID"})
	}

	action, err := h.actionService.GetActionByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Action not found"})
	}

	response := dto.ToActionResponse(action)
	return c.JSON(http.StatusOK, response)
}

// GetActionsByResourceID retrieves actions by resource ID
// @Summary Get actions by resource ID
// @Description Retrieve all actions associated with a specific resource
// @Tags actions
// @Accept json
// @Produce json
// @Param resourceID path string true "Resource ID"
// @Success 200 {array} dto.ActionResponse "Actions found"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "No actions found"
// @Router /api/actions/resource/{resourceID} [get]
func (h *ActionHandler) GetActionsByResourceID(c echo.Context) error {
	resourceID := c.Param("resourceID")
	if resourceID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing resource ID"})
	}

	actions, err := h.actionService.GetActionsByResourceID(c.Request().Context(), resourceID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if len(actions) == 0 {
		return c.JSON(http.StatusOK, []dto.ActionResponse{})
	}

	// Convert domain models to response DTOs
	actionResponses := make([]dto.ActionResponse, len(actions))
	for i, a := range actions {
		actionResponses[i] = dto.ToActionResponse(a)
	}

	return c.JSON(http.StatusOK, actionResponses)
}

// ListActions retrieves a paginated list of actions
// @Summary List actions
// @Description Get a paginated list of all actions
// @Tags actions
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to return (default: 10)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Success 200 {object} dto.ListActionsResponse "List of actions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/actions [get]
func (h *ActionHandler) ListActions(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	actions, err := h.actionService.ListActions(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert domain models to response DTOs
	actionResponses := make([]dto.ActionResponse, len(actions))
	for i, a := range actions {
		actionResponses[i] = dto.ToActionResponse(a)
	}

	response := dto.ListActionsResponse{
		Actions: actionResponses,
		Total:   len(actionResponses), // In a real app, we'd get the total count from the service
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateAction updates an existing action
// @Summary Update an action
// @Description Update an existing action by its ID
// @Tags actions
// @Accept json
// @Produce json
// @Param id path string true "Action ID"
// @Param action body dto.UpdateActionRequest true "Updated action information"
// @Success 200 {object} dto.ActionResponse "Action updated"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Action not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/actions/{id} [put]
func (h *ActionHandler) UpdateAction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing action ID"})
	}

	var req dto.UpdateActionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get existing action
	action, err := h.actionService.GetActionByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Action not found"})
	}

	// Update the action with request data
	req.UpdateActionDomain(action)

	if err := h.actionService.UpdateAction(c.Request().Context(), action); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToActionResponse(action)
	return c.JSON(http.StatusOK, response)
}

// DeleteAction deletes an action by ID
// @Summary Delete an action
// @Description Delete an action by its ID
// @Tags actions
// @Accept json
// @Produce json
// @Param id path string true "Action ID"
// @Success 200 {object} dto.ActionResponse "Action soft deleted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Action not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/actions/{id} [delete]
func (h *ActionHandler) DeleteAction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing action ID"})
	}

	deletedAction, err := h.actionService.DeleteAction(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "action not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Action not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToActionResponse(deletedAction)
	return c.JSON(http.StatusOK, response)
}
