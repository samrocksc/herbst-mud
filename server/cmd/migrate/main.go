package main

// cmd/migrate is a standalone schema-migration entrypoint. It connects to
// PostgreSQL using the same environment variables as the web server and runs
// the ent schema migration (client.Schema.Create) against the target database.
//
// It is intended to run as a one-shot container/service in the Docker Compose
// stack, before the API server starts, so schema changes apply deterministically
// on every deployment (dev and prod).
//
// Usage:
//
//	cmd/migrate
//
// Configuration (same as the web server):
//
//	DATABASE_URL  — single Postgres DSN (takes precedence)
//	DB_HOST DB_PORT DB_USER DB_PASSWORD DB_NAME DB_SSL_MODE — component vars
//
// Without any env vars set it defaults to a local dev database
// (localhost:5432, herbst/herbst_password/herbst_mud).

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // register the postgres driver

	"herbst-server/db"
)

func getDBConfig() string {
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	isDev := host == "" && user == ""
	if isDev {
		log.Println("Warning: Using development database defaults. Set DATABASE_URL or DB_* env vars for production.")
		host = "localhost"
		port = "5432"
		user = "herbst"
		password = "herbst_password"
		dbname = "herbst_mud"
	}

	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		if isDev {
			sslMode = "disable"
		} else {
			sslMode = "require"
		}
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode)
}

func main() {
	client, err := db.Open("postgres", getDBConfig())
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("Database schema migrated successfully")
}
