package questservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CheckProgress checks if a game event advances any of the character's active quests.
// It calls the server's bulk check endpoint which finds matching quests and
// increments progress on each one.
func (s *Service) CheckProgress(charID int, objectiveType string, targetID string) error {
	url := fmt.Sprintf("%s/api/characters/%d/quests/check-all", s.restBase, charID)
	body := map[string]interface{}{
		"objective_type": objectiveType,
		"target_id":      targetID,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := s.client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("POST %s: %d", url, resp.StatusCode)
	}
	return nil
}

// AcceptQuest accepts a quest for a character via the REST API.
func (s *Service) AcceptQuest(charID int, questID int) (*QuestProgress, error) {
	url := fmt.Sprintf("%s/api/characters/%d/quests", s.restBase, charID)
	body := map[string]interface{}{"quest_id": questID}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("POST %s: %d", url, resp.StatusCode)
	}
	var qp QuestProgress
	if err := json.NewDecoder(resp.Body).Decode(&qp); err != nil {
		return nil, err
	}
	return &qp, nil
}

// AbandonQuest abandons a quest for a character via the REST API.
func (s *Service) AbandonQuest(charID int, questID int) error {
	url := fmt.Sprintf("%s/api/characters/%d/quests/%d/abandon", s.restBase, charID, questID)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PUT %s: %d", url, resp.StatusCode)
	}
	return nil
}