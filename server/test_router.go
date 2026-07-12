package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"herbst-server/db"
	"herbst-server/middleware"
	"herbst-server/repository"
	"herbst-server/routes"
	"herbst-server/service"
)

// setupTestRouter builds a Gin router with the real route registrations used by
// the application. Tests should call this instead of hand-registering a subset
// of routes, which drifts out of sync with main.go.
func setupTestRouter(client *db.Client) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Ensure a deterministic JWT secret for token validation, but respect an
	// explicitly configured secret so individual tests can set their own.
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-secret-key-for-integration-tests")
	}

	router := gin.New()

	repos := repository.NewContainer(client)
	services := service.NewContainer(client, repos, slog.Default())

	// Public / auth-only routes.
	routes.RegisterUserRoutes(router, repos)

	// Character routes (mix of public playable endpoints and protected CRUD).
	routes.RegisterCharacterRoutes(router, services, repos, client)

	// Protected /api routes that need a valid JWT.
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(nil))
	{
		routes.RegisterMeRoutes(protected, services, repos)
	}

	// Admin-only /api routes.
	routes.RegisterRoomRoutes(router, client, services)
	routes.RegisterItemInstanceRoutes(router, repos, client)

	// Other real route groups used by integration tests.
	routes.RegisterEquipmentRoutes(router, repos, client)
	routes.RegisterEquipmentEquipRoutes(router, repos)
	routes.RegisterAbilityRoutes(router, repos, client)
	routes.RegisterEffectRoutes(router, repos, client)
	routes.RegisterAdminWipeRoutes(router, client)

	return router
}

// adminToken returns a signed JWT with the admin flag set. Tests that call
// admin-only routes can use this instead of hand-crafting claims.
func adminToken() string {
	claims := &middleware.Claims{
		UserID:  1,
		Email:   "admin@example.com",
		IsAdmin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test-secret-key-for-integration-tests"
	}
	s, _ := token.SignedString([]byte(secret))
	return s
}
