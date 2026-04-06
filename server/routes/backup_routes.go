package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/backup"
	"herbst-server/db"
)

// RegisterBackupRoutes registers backup-related API endpoints
func RegisterBackupRoutes(router *gin.Engine, client *db.Client) {
	backupDir := "./backups"
	os.MkdirAll(backupDir, 0755)

	router.POST("/api/backups", createBackupHandler(client, backupDir))
	router.GET("/api/backups", listBackupsHandler(backupDir))
	router.GET("/api/backups/:id/validate", validateBackupHandler(backupDir))
	router.POST("/api/backups/:id/restore", restoreBackupHandler(client, backupDir))
	router.GET("/api/backups/:id/manifest", getManifestHandler(backupDir))
}

func createBackupHandler(client *db.Client, backupDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := backup.CreateBackup(client, backupDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message":  "Backup created successfully",
			"path":     result.Path,
			"manifest": result.Manifest,
		})
	}
}

func listBackupsHandler(backupDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		backups, err := backup.ListBackups(backupDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"backups": backups})
	}
}

func validateBackupHandler(backupDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		backupID := c.Param("id")
		if strings.Contains(backupID, "..") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
			return
		}

		backupPath := filepath.Join(backupDir, backupID)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
			return
		}

		result, err := backup.ValidateBackup(backupPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}