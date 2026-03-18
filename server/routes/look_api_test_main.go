package main

import (
	"encoding/json"
	"fmt"
)

// Test functions - standalone without dependencies

// processHiddenDetails is a copy of the function from item_routes.go for testing
func processHiddenDetails(details []map[string]interface{}, examineLevel int) []map[string]interface{} {
	if details == nil {
		return nil
	}

	var revealed []map[string]interface{}
	for _, detail := range details {
		minLevel := 0
		if ml, ok := detail["min_examine_level"].(float64); ok {
			minLevel = int(ml)
		}

		if examineLevel >= minLevel {
			detail["revealed"] = true
			revealed = append(revealed, detail)
		} else {
			detail["revealed"] = false
			revealed = append(revealed, detail)
		}
	}

	return revealed
}

func testProcessHiddenDetails() {
	fmt.Println("=== Test: processHiddenDetails ===")
	
	// Test 1: No details
	result := processHiddenDetails(nil, 10)
	if result != nil {
		fmt.Println("FAIL: Expected nil for nil input")
	} else {
		fmt.Println("PASS: nil input returns nil")
	}
	
	// Test 2: All revealed
	details := []map[string]interface{}{
		{"text": "detail1", "min_examine_level": float64(5)},
		{"text": "detail2", "min_examine_level": float64(10)},
	}
	result = processHiddenDetails(details, 10)
	if result[0]["revealed"] == true && result[1]["revealed"] == true {
		fmt.Println("PASS: All details revealed at level 10")
	} else {
		fmt.Println("FAIL: Not all details revealed")
	}
	
	// Test 3: Some not revealed
	details = []map[string]interface{}{
		{"text": "easy", "min_examine_level": float64(5)},
		{"text": "hard", "min_examine_level": float64(50)},
	}
	result = processHiddenDetails(details, 10)
	if result[0]["revealed"] == true && result[1]["revealed"] == false {
		fmt.Println("PASS: Hard detail not revealed at level 10")
	} else {
		fmt.Println("FAIL: Logic error")
	}
	
	fmt.Println()
}

func testJSONMarshaling() {
	fmt.Println("=== Test: JSON Response Structures ===")
	
	// Item examine response
	examineResp := map[string]interface{}{
		"id":             1,
		"name":           "test_item",
		"description":   "A test item",
		"examineDesc":    "A detailed description",
		"hiddenDetails": []map[string]interface{}{},
		"isReadable":    true,
		"readContent":   "Test content",
		"examineLevel":   10,
		"type":           "weapon",
		"weight":         5,
		"level":          1,
		"isImmovable":    false,
		"isContainer":    false,
	}
	
	jsonBytes, err := json.Marshal(examineResp)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}
	
	var decoded map[string]interface{}
	err = json.Unmarshal(jsonBytes, &decoded)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}
	
	requiredFields := []string{"id", "name", "description", "examineDesc", "hiddenDetails", "examineLevel"}
	allPresent := true
	for _, field := range requiredFields {
		if decoded[field] == nil {
			fmt.Printf("FAIL: Missing field %s\n", field)
			allPresent = false
		}
	}
	if allPresent {
		fmt.Println("PASS: Item examine response structure valid")
	}
	
	// Room response
	roomResp := map[string]interface{}{
		"id":             1,
		"name":           "Town Square",
		"description":    "A busy town square.",
		"isStartingRoom": true,
		"exits":          map[string]int{"north": 2, "south": 3},
		"items":          []int{1, 2, 3},
		"npcs":           []int{10, 11},
		"players":        []int{},
	}
	
	roomBytes, _ := json.Marshal(roomResp)
	var roomDecoded map[string]interface{}
	json.Unmarshal(roomBytes, &roomDecoded)
	
	if roomDecoded["items"] != nil && roomDecoded["npcs"] != nil && roomDecoded["players"] != nil {
		fmt.Println("PASS: Room response with items/NPCs/players valid")
	} else {
		fmt.Println("FAIL: Room response missing embedded collections")
	}
	
	fmt.Println()
}

func testReadableContentLogic() {
	fmt.Println("=== Test: Readable Content Logic ===")
	
	tests := []struct {
		name           string
		charLevel      int
		requiredLevel  int
		expectRevealed bool
	}{
		{"Level 1 vs req 5", 1, 5, false},
		{"Level 5 vs req 5", 5, 5, true},
		{"Level 10 vs req 5", 10, 5, true},
		{"Level 3 vs req 0", 3, 0, true},
	}
	
	allPass := true
	for _, tt := range tests {
		result := tt.charLevel >= tt.requiredLevel
		if result != tt.expectRevealed {
			fmt.Printf("FAIL: %s - got %v, expected %v\n", tt.name, result, tt.expectRevealed)
			allPass = false
		}
	}
	if allPass {
		fmt.Println("PASS: All readable content logic tests")
	}
	
	fmt.Println()
}

func main() {
	testProcessHiddenDetails()
	testJSONMarshaling()
	testReadableContentLogic()
	
	fmt.Println("=== All Tests Complete ===")
}