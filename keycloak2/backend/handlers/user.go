package handlers

import (
	"net/http"

	"keycloak-tutorial-backend/middleware"
	"keycloak-tutorial-backend/models"

	"github.com/labstack/echo/v4"
)

// GetUserInfo returns authenticated user information
func GetUserInfo(c echo.Context) error {
	// Get user claims from context (set by JWT middleware)
	userClaims, ok := c.Get("user").(*middleware.KeycloakClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not found in context")
	}

	// Create response with user information
	userResponse := models.UserResponse{
		Name:  userClaims.Name,
		Email: userClaims.Email,
	}

	return c.JSON(http.StatusOK, userResponse)
}
