package routes

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"herbst-server/backup"
	"herbst-server/db"
	"herbst-server/dblog"
)

func restoreBackupHandler(client *db.Client, backupDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		backupID := c.Param("id")
		if strings.Contains(backupID, "..") {
			slog.Warn("bad request", slog.String("service", "backup"), slog.String("reason", "Invalid backup ID"), slog.String("client_ip", c.ClientIP()))
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
			dblog.Error("backup validation failed", err, slog.String("service", "backup"), slog.String("backup_id", backupID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !result.Valid {
			slog.Warn("bad request", slog.String("service", "backup"), slog.String("reason", "Backup validation failed"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Backup validation failed",
				"details": result,
			})
			return
		}

		if err := backup.RestoreBackup(client, backupPath); err != nil {
			dblog.Error("backup restore failed", err, slog.String("service", "backup"), slog.String("backup_id", backupID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		slog.Info("backup restored", slog.String("backup_id", backupID), slog.String("user_email", c.GetString("email")), slog.String("service", "backup"))
		c.JSON(http.StatusOK, gin.H{"message": "Backup restored successfully"})
	}
}

func getManifestHandler(backupDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		backupID := c.Param("id")
		if strings.Contains(backupID, "..") {
			slog.Warn("bad request", slog.String("service", "backup"), slog.String("reason", "Invalid backup ID"), slog.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
			return
		}

		backupPath := filepath.Join(backupDir, backupID)
		manifestPath := filepath.Join(backupPath, "manifest.json")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Backup manifest not found"})
			return
		}

		c.Data(http.StatusOK, "application/json", data)
	}
}