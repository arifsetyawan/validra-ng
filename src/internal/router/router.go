package router

import (
	"net/http"
	"time"

	"github.com/arifsetyawan/validra/src/internal/delivery/http/handler"
	"github.com/arifsetyawan/validra/src/internal/service"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Register registers all routes and handlers to the echo instance
func Register(e *echo.Echo, resourceService *service.ResourceService, userService *service.UserService, roleService *service.RoleService, actionService *service.ActionService) {
	// API routes
	registerAPIRoutes(e, resourceService, userService, roleService, actionService)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Swagger documentation routes
	registerSwaggerRoutes(e)
}

// registerAPIRoutes sets up all API-related routes
func registerAPIRoutes(e *echo.Echo, resourceService *service.ResourceService, userService *service.UserService, roleService *service.RoleService, actionService *service.ActionService) {
	// Initialize handlers
	resourceHandler := handler.NewResourceHandler(resourceService)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	actionHandler := handler.NewActionHandler(actionService)

	// Register routes for each handler
	resourceHandler.Register(e)
	userHandler.Register(e)
	roleHandler.Register(e)
	actionHandler.Register(e)
}

// registerSwaggerRoutes sets up Swagger documentation routes
func registerSwaggerRoutes(e *echo.Echo) {
	e.GET("/docs/*", echoSwagger.WrapHandler)
}
