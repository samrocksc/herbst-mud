package validate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"herbst-server/backup/types"
)

func validateCharacters(backupDir string, userIDs, roomIDs, characterIDs map[int]bool, errors *[]types.ValidationError) {
	data, err := os.ReadFile(filepath.Join(backupDir, "characters.json"))
	if err != nil {
		return
	}

	var characters []map[string]interface{}
	if err := json.Unmarshal(data, &characters); err != nil {
		return
	}

	for _, c := range characters {
		if id, ok := c["id"].(float64); ok {
			characterIDs[int(id)] = true
		}
		checkCharacterRefs(c, userIDs, roomIDs, errors)
	}
}

func checkCharacterRefs(c map[string]interface{}, userIDs, roomIDs map[int]bool, errors *[]types.ValidationError) {
	if userID, ok := c["user_id"].(float64); ok && int(userID) != 0 {
		if !userIDs[int(userID)] {
			*errors = append(*errors, types.ValidationError{
				File:     "characters.json",
				Severity: "warning",
				Message:  fmt.Sprintf("Character references non-existent user_id: %d", int(userID)),
			})
		}
	}
	if roomID, ok := c["current_room_id"].(float64); ok && int(roomID) != 0 {
		if !roomIDs[int(roomID)] {
			*errors = append(*errors, types.ValidationError{
				File:     "characters.json",
				Severity: "warning",
				Message:  fmt.Sprintf("Character references non-existent current_room_id: %d", int(roomID)),
			})
		}
	}
}

func validateEquipment(backupDir string, roomIDs map[int]bool, errors *[]types.ValidationError) {
	data, err := os.ReadFile(filepath.Join(backupDir, "equipment.json"))
	if err != nil {
		return
	}

	var equipment []map[string]interface{}
	if err := json.Unmarshal(data, &equipment); err != nil {
		return
	}

	for _, e := range equipment {
		if roomID, ok := e["room_id"].(float64); ok && int(roomID) != 0 {
			if !roomIDs[int(roomID)] {
				*errors = append(*errors, types.ValidationError{
					File:     "equipment.json",
					Severity: "warning",
					Message:  fmt.Sprintf("Equipment references non-existent room_id: %d", int(roomID)),
				})
			}
		}
	}
}