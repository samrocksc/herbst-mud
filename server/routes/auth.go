package routes

import "os"

// getJWTSecret returns the JWT secret from environment or uses default
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// In production, this should be set via environment variable
		return "herbst-mud-secret-key-change-in-production"
	}
	return secret
}