package middleware

import (
	"github.com/arifsetyawan/validra/src/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware configures the middleware for the Echo instance
func SetupMiddleware(e *echo.Echo, logger *logger.Logger) {
	// Recover from panics
	e.Use(middleware.Recover())

	// Request ID
	e.Use(middleware.RequestID())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Logging middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
	}))

	// Custom context middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// You can add custom logic here if needed
			return next(c)
		}
	})
}
