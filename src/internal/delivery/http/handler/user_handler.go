package handler

import (
	"net/http"
	"strconv"

	"github.com/arifsetyawan/validra/src/internal/delivery/http/dto"
	"github.com/arifsetyawan/validra/src/internal/service"
	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register registers the routes to the given echo instance
func (h *UserHandler) Register(e *echo.Echo) {
	users := e.Group("/api/users")
	users.POST("", h.CreateUser)
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User information"
// @Success 201 {object} dto.UserResponse "User created"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user := req.ToUserDomain()
	if err := h.userService.CreateUser(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToUserResponse(user)
	return c.JSON(http.StatusCreated, response)
}

// GetUser retrieves a user by ID
// @Summary Get a user by ID
// @Description Retrieve a specific user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse "User found"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing user ID"})
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	response := dto.ToUserResponse(user)
	return c.JSON(http.StatusOK, response)
}

// ListUsers retrieves a paginated list of users
// @Summary List users
// @Description Get a paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Number of items to return (default: 10)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Success 200 {object} dto.ListUsersResponse "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users [get]
func (h *UserHandler) ListUsers(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	users, err := h.userService.ListUsers(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert domain models to response DTOs
	userResponses := make([]dto.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = dto.ToUserResponse(u)
	}

	response := dto.ListUsersResponse{
		Users: userResponses,
		Total: len(userResponses), // In a real app, we'd get the total count from the service
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateUser updates an existing user
// @Summary Update a user
// @Description Update an existing user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "Updated user information"
// @Success 200 {object} dto.UserResponse "User updated"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing user ID"})
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Get existing user
	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Update the user with request data
	req.UpdateUserDomain(user)

	if err := h.userService.UpdateUser(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := dto.ToUserResponse(user)
	return c.JSON(http.StatusOK, response)
}

// DeleteUser deletes a user by ID
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "User deleted"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing user ID"})
	}

	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
