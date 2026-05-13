package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"herbst-server/middleware"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: token <user_id> <email> <is_admin>")
		fmt.Println("  user_id:  Integer user/character ID")
		fmt.Println("  email:    Email address (e.g., test@example.com)")
		fmt.Println("  is_admin: true or false")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("  token 1 sma@example.com true")
		os.Exit(1)
	}

	userID := os.Args[1]
	email := os.Args[2]
	isAdminStr := os.Args[3]

	// Parse is_admin
	isAdmin := isAdminStr == "true" || isAdminStr == "1" || isAdminStr == "yes"

	// Get JWT secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("Warning: JWT_SECRET not set, using default development secret")
		secret = "dev-secret-key-not-for-production-use-only"
	}

	// Create token
	claims := &middleware.Claims{
		UserID:  1, // Will be parsed from args
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// Parse user ID
	var uid uint
	fmt.Sscanf(userID, "%d", &uid)
	claims.UserID = uid

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	// Output in format suitable for curl or Bearer header
	fmt.Println("")
	fmt.Println("=== Bearer Token (for API debugging) ===")
	fmt.Println("")
	fmt.Printf("Token: %s\n", tokenString)
	fmt.Println("")
	fmt.Println("Usage with curl:")
	fmt.Printf("  curl -H 'Authorization: Bearer %s' http://localhost:8080/api/...\n", tokenString)
	fmt.Println("")
	fmt.Println("Usage in herbst SSH client (set TOKEN env var):")
	fmt.Printf("  export TOKEN='%s'\n", tokenString)
	fmt.Println("  Then restart the herbst client to use it")
	fmt.Println("")
}
