package handler

import (
	"net/http"
	"strconv"

	"github.com/arifsetyawan/validra/src/internal/dto"
	"github.com/arifsetyawan/validra/src/internal/service"
	"github.com/labstack/echo/v4"
)

// ResourceHandler handles HTTP requests for resources
type ResourceHandler struct {
	resourceService *service.ResourceService
}

// NewResourceHandler creates a new ResourceHandler
func NewResourceHandler(resourceService *service.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

// Register registers the routes to the given echo instance
func (h *ResourceHandler) Register(e *echo.Echo) {
	resources := e.Group("/api/resources")
	resources.POST("", h.CreateResource)
	resources.GET("", h.ListResources)
	resources.GET("/:id", h.GetResource)
	resources.PUT("/:id", h.UpdateResource)
	resources.DELETE("/:id", h.DeleteResource)
}

// CreateResource creates a new resource
// @Summary Create a new resource
// @Description Create a new resource with the provided information
// @Tags resources
// @Accept json
// @Produce json
// @Param resource body dto.CreateResourceRequest true "Resource information"
// @Success 201 {object} dto.ResourceResponse "Resource created"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/resources [post]
func (h *ResourceHandler) CreateResource(c echo.Context) error {
	var req dto.CreateResourceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	resource := req.ToResourceDomain()
	if err := h.resourceService.CreateResource(c.Request().Context(), resource); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToResourceResponse(resource)
	return c.JSON(http.StatusCreated, response)
}

// GetResource retrieves a resource by ID
// @Summary Get a resource by ID
// @Description Retrieve a specific resource by its unique identifier
// @Tags resources
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Success 200 {object} dto.ResourceResponse "Resource found"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Resource not found"
// @Router /api/resources/{id} [get]
func (h *ResourceHandler) GetResource(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing resource ID"})
	}

	resource, err := h.resourceService.GetResourceByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Resource not found"})
	}

	response := dto.ToResourceResponse(resource)
	return c.JSON(http.StatusOK, response)
}

// ListResources retrieves a paginated list of resources
// @Summary List resources
// @Description Get a paginated list of all resources
// @Tags resources
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to return (default: 10)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Success 200 {object} dto.ListResourcesResponse "List of resources"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/resources [get]
func (h *ResourceHandler) ListResources(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	resources, err := h.resourceService.ListResources(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert domain models to response DTOs
	resourceResponses := make([]dto.ResourceResponse, len(resources))
	for i, r := range resources {
		resourceResponses[i] = dto.ToResourceResponse(r)
	}

	response := dto.ListResourcesResponse{
		Resources: resourceResponses,
		Total:     len(resourceResponses), // In a real app, we'd get the total count from the service
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateResource updates an existing resource
// @Summary Update a resource
// @Description Update an existing resource by its ID
// @Tags resources
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Param resource body dto.UpdateResourceRequest true "Updated resource information"
// @Success 200 {object} dto.ResourceResponse "Resource updated"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Resource not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/resources/{id} [put]
func (h *ResourceHandler) UpdateResource(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing resource ID"})
	}

	var req dto.UpdateResourceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get existing resource
	resource, err := h.resourceService.GetResourceByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Resource not found"})
	}

	// Update the resource with request data
	req.UpdateResourceDomain(resource)

	if err := h.resourceService.UpdateResource(c.Request().Context(), resource); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToResourceResponse(resource)
	return c.JSON(http.StatusOK, response)
}

// DeleteResource deletes a resource by ID
// @Summary Delete a resource
// @Description Delete a resource by its ID
// @Tags resources
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Success 200 {object} dto.ResourceResponse "Resource soft deleted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Resource not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/resources/{id} [delete]
func (h *ResourceHandler) DeleteResource(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing resource ID"})
	}

	deletedResource, err := h.resourceService.DeleteResource(c.Request().Context(), id)
	if err != nil {
		if err.Error() == "resource not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Resource not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToResourceResponse(deletedResource)
	return c.JSON(http.StatusOK, response)
}
