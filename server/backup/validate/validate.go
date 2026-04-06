package validate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"herbst-server/backup/types"
)

// Backup validates a backup directory for integrity and correctness
func Backup(backupDir string) (*types.ValidationResult, error) {
	result := &types.ValidationResult{
		Valid:    true,
		Errors:   []types.ValidationError{},
		Warnings: []types.ValidationError{},
	}

	manifest, err := parseManifest(backupDir)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, types.ValidationError{
			File:     "manifest.json",
			Severity: "critical",
			Message:  fmt.Sprintf("Failed to parse manifest: %v", err),
		})
		return result, nil
	}

	for entity, fileName := range types.EntityFileNames {
		filePath := filepath.Join(backupDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			result.Valid = false
			result.Errors = append(result.Errors, types.ValidationError{
				File:     fileName,
				Severity: "critical",
				Message:  fmt.Sprintf("Missing required file for entity: %s", entity),
			})
		}
	}

	checksumErrors := verifyChecksums(backupDir, manifest)
	for _, ce := range checksumErrors {
		result.Warnings = append(result.Warnings, ce)
	}

	integrityErrors := referentialIntegrity(backupDir)
	for _, ie := range integrityErrors {
		if ie.Severity == "critical" {
			result.Valid = false
			result.Errors = append(result.Errors, ie)
		} else {
			result.Warnings = append(result.Warnings, ie)
		}
	}

	return result, nil
}

func parseManifest(backupDir string) (*types.Manifest, error) {
	manifestPath := filepath.Join(backupDir, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var manifest types.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func verifyChecksums(backupDir string, manifest *types.Manifest) []types.ValidationError {
	var errors []types.ValidationError

	for fileName, expectedChecksum := range manifest.Checksums {
		filePath := filepath.Join(backupDir, fileName)
		actualChecksum, err := calculateChecksum(filePath)
		if err != nil {
			errors = append(errors, types.ValidationError{
				File:     fileName,
				Severity: "warning",
				Message:  fmt.Sprintf("Failed to calculate checksum: %v", err),
			})
			continue
		}

		if actualChecksum != expectedChecksum {
			errors = append(errors, types.ValidationError{
				File:     fileName,
				Severity: "warning",
				Message:  fmt.Sprintf("Checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum),
			})
		}
	}

	return errors
}