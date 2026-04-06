package backup

import (
	"herbst-server/backup/types"
	"herbst-server/backup/validate"
)

// ValidateBackup validates a backup directory for integrity and correctness
func ValidateBackup(backupDir string) (*types.ValidationResult, error) {
	return validate.Backup(backupDir)
}