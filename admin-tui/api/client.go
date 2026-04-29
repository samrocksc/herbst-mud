package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Token file path for persistent JWT storage
var tokenFile = os.Getenv("HOME") + "/.config/herbst-admin/token"

// GetToken reads the stored JWT token from disk
func GetToken() string {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// SetToken stores the JWT token to disk
func SetToken(token string) error {
	// Ensure directory exists
	dir := filepath.Dir(tokenFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(tokenFile, []byte(token), 0600)
}

// BaseURL is the REST API server URL.
// Set by the main package on startup.
var BaseURL = "http://localhost:8080"

var httpClient = &http.Client{Timeout: 10 * time.Second}

// AuthPayload is the login request body
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the login response
type AuthResponse struct {
	Token   string `json:"token"`
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}

// User represents a user account
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
	CreatedAt string `json:"created_at,omitempty"`
}

// Character represents a game character
type Character struct {
	ID          int    `json:"id"`
	Name       string `json:"name"`
	Level      int    `json:"level"`
	Class      string `json:"class"`
	Race       string `json:"race"`
	HP         int    `json:"hp"`
	MaxHP      int    `json:"max_hp"`
	RoomID     int    `json:"room_id"`
	OwnerID    int    `json:"owner_id"`
	IsNPC      bool   `json:"is_npc"`
	Description string `json:"description"`
	Behavior   string `json:"behavior"`
	Aggression string `json:"aggression"`
}

// Room represents a game room
type Room struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Exits       map[string]int `json:"exits"`
	IsStarting  bool           `json:"isStartingRoom"`
	Floor       int            `json:"floor,omitempty"`
	Edges       map[string]any `json:"edges,omitempty"`
}

// EquipmentItem represents an item in the world
type EquipmentItem struct {
	ID          int    `json:"id"`
	Name       string `json:"name"`
	Description string `json:"description"`
	Slot       string `json:"slot"`
	ItemType   string `json:"itemType"`
	Level      int    `json:"level"`
	Weight     int    `json:"weight"`
	Color      string `json:"color"`
	RoomID     int    `json:"room_id,omitempty"`
	OwnerID    *int   `json:"ownerId,omitempty"`
	IsVisible  bool   `json:"isVisible"`
	IsEquipped bool   `json:"isEquipped"`
	Healing    int    `json:"healing"`
	Effect     string `json:"effect,omitempty"`
}

// Quest represents a quest
type Quest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	 Objectives  []string `json:"objectives,omitempty"`
	Giver       string   `json:"giver,omitempty"`
	Rewards     string   `json:"rewards,omitempty"`
}

// Talent represents a combat talent
type Talent struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Cooldown    int      `json:"cooldown"`
	ManaCost    int      `json:"manaCost"`
	StaminaCost int      `json:"staminaCost"`
	EffectType  string   `json:"effectType"`
	Tags        []string `json:"tags,omitempty"`
}

// HealthResponse is the server health check response
type HealthResponse struct {
	Status string `json:"status"`
	SSH    string `json:"ssh"`
	DB     string `json:"db"`
}

// BackupManifest describes a backup
type BackupManifest struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	CreatedAt string `json:"created_at"`
	Size      int64  `json:"size"`
}

// FactionCategory represents a faction category
type FactionCategory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// Faction represents a player/NPC faction
type Faction struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	CategoryID  int      `json:"category_id,omitempty"`
	Standing    int      `json:"standing"`
	Members     []int    `json:"members,omitempty"`
	IsUniversal bool     `json:"is_universal,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
}

// CharacterFaction represents a character's faction membership
type CharacterFaction struct {
	ID         int    `json:"id"`
	CharacterID int    `json:"character_id"`
	FactionID  int    `json:"faction_id"`
	FactionName string `json:"faction_name,omitempty"`
	Standing   int    `json:"standing"`
	JoinedAt   string `json:"joined_at,omitempty"`
}

// APIError wraps server error responses
type APIError struct {
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func readResponse[T any](resp *http.Response) (T, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *new(T), fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return *new(T), APIError{StatusCode: resp.StatusCode, Message: string(body)}
	}
	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return *new(T), fmt.Errorf("unmarshal: %w", err)
	}
	return result, nil
}

func doRequest[T any](method, path string, body any) (T, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return *new(T), fmt.Errorf("marshal: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, BaseURL+path, reqBody)
	if err != nil {
		return *new(T), fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token := GetToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return *new(T), fmt.Errorf("do request: %w", err)
	}
	return readResponse[T](resp)
}

// ─── Auth ───────────────────────────────────────────────────────────────────

// Login authenticates and stores the JWT token locally
func Login(email, password string) (*AuthResponse, error) {
	resp, err := httpClient.Post(BaseURL+"/users/auth", "application/json",
		bytes.NewReader(must(json.Marshal(AuthPayload{Email: email, Password: password}))))
	if err != nil {
		return nil, fmt.Errorf("login request: %w", err)
	}
	result, err := readResponse[AuthResponse](resp)
	if err != nil {
		return nil, err
	}
	SetToken(result.Token)
	return &result, nil
}

// ─── Health ─────────────────────────────────────────────────────────────────

func Health() (HealthResponse, error) {
	return doRequest[HealthResponse]("GET", "/healthz", nil)
}

// ─── Users ─────────────────────────────────────────────────────────────────

func ListUsers() ([]User, error) {
	return doRequest[[]User]("GET", "/users", nil)
}

func UpdateUser(id int, body map[string]any) (User, error) {
	return doRequest[User]("PUT", fmt.Sprintf("/users/%d", id), body)
}

// ─── Characters ─────────────────────────────────────────────────────────────

func ListCharacters() ([]Character, error) {
	return doRequest[[]Character]("GET", "/characters", nil)
}

func GetCharacter(id int) (Character, error) {
	return doRequest[Character]("GET", fmt.Sprintf("/characters/%d", id), nil)
}

func UpdateCharacter(id int, body map[string]any) (Character, error) {
	return doRequest[Character]("PUT", fmt.Sprintf("/characters/%d", id), body)
}

func DeleteCharacter(id int) error {
	_, err := doRequest[any]("DELETE", fmt.Sprintf("/characters/%d", id), nil)
	return err
}

// ─── Rooms ──────────────────────────────────────────────────────────────────

func ListRooms() ([]Room, error) {
	return doRequest[[]Room]("GET", "/rooms", nil)
}

func GetRoom(id int) (Room, error) {
	return doRequest[Room]("GET", fmt.Sprintf("/rooms/%d", id), nil)
}

func CreateRoom(body map[string]any) (Room, error) {
	return doRequest[Room]("POST", "/rooms", body)
}

func UpdateRoom(id int, body map[string]any) (Room, error) {
	return doRequest[Room]("PUT", fmt.Sprintf("/rooms/%d", id), body)
}

func DeleteRoom(id int) error {
	_, err := doRequest[any]("DELETE", fmt.Sprintf("/rooms/%d", id), nil)
	return err
}

// ─── NPCs ───────────────────────────────────────────────────────────────────

func ListNPCs() ([]Character, error) {
	chars, err := doRequest[[]Character]("GET", "/npcs", nil)
	return chars, err
}

func CreateNPC(body map[string]any) (Character, error) {
	return doRequest[Character]("POST", "/npcs", body)
}

func UpdateNPC(id int, body map[string]any) (Character, error) {
	return doRequest[Character]("PUT", fmt.Sprintf("/characters/%d", id), body)
}

func DeleteNPC(id int) error {
	_, err := doRequest[any]("DELETE", fmt.Sprintf("/characters/%d", id), nil)
	return err
}

// ─── Equipment / Items ────────────────────────────────────────────────────────

func ListItems() ([]EquipmentItem, error) {
	return doRequest[[]EquipmentItem]("GET", "/equipment", nil)
}

func CreateItem(body map[string]any) (EquipmentItem, error) {
	return doRequest[EquipmentItem]("POST", "/equipment", body)
}

func UpdateItem(id int, body map[string]any) (EquipmentItem, error) {
	return doRequest[EquipmentItem]("PUT", fmt.Sprintf("/equipment/%d", id), body)
}

func DeleteItem(id int) error {
	_, err := doRequest[any]("DELETE", fmt.Sprintf("/equipment/%d", id), nil)
	return err
}

// ─── Quests ─────────────────────────────────────────────────────────────────

func ListQuests() ([]Quest, error) {
	return doRequest[[]Quest]("GET", "/content/quests", nil)
}

// ─── Skills & Talents ────────────────────────────────────────────────────────

func ListTalents() ([]Talent, error) {
	return doRequest[[]Talent]("GET", "/talents", nil)
}

// ─── Backup ─────────────────────────────────────────────────────────────────

func ListBackups() ([]BackupManifest, error) {
	return doRequest[[]BackupManifest]("GET", "/api/backups", nil)
}

func TriggerBackup() error {
	_, err := doRequest[any]("POST", "/api/backups", nil)
	return err
}

// ─── Factions ────────────────────────────────────────────────────────────────

func ListFactions() ([]Faction, error) {
	return doRequest[[]Faction]("GET", "/api/factions", nil)
}

func GetFaction(id int) (Faction, error) {
	return doRequest[Faction]("GET", fmt.Sprintf("/api/factions/%d", id), nil)
}

func CreateFaction(body map[string]any) (Faction, error) {
	return doRequest[Faction]("POST", "/api/factions", body)
}

func UpdateFaction(id int, body map[string]any) (Faction, error) {
	return doRequest[Faction]("PUT", fmt.Sprintf("/api/factions/%d", id), body)
}

func DeleteFaction(id int) error {
	_, err := doRequest[any]("DELETE", fmt.Sprintf("/api/factions/%d", id), nil)
	return err
}

func AssignCharacterToFaction(characterID, factionID int) error {
	_, err := doRequest[any]("POST", fmt.Sprintf("/api/factions/%d/assign", factionID),
		map[string]any{"character_id": characterID})
	return err
}

func RemoveCharacterFromFaction(characterID, factionID int) error {
	_, err := doRequest[any]("POST", fmt.Sprintf("/api/factions/%d/unassign", factionID),
		map[string]any{"character_id": characterID})
	return err
}

// ─── Faction Categories ─────────────────────────────────────────────────────

func ListFactionCategories() ([]FactionCategory, error) {
	return doRequest[[]FactionCategory]("GET", "/api/faction-categories", nil)
}

func CreateFactionCategory(body map[string]any) (FactionCategory, error) {
	return doRequest[FactionCategory]("POST", "/api/faction-categories", body)
}

func UpdateFactionCategory(id int, body map[string]any) (FactionCategory, error) {
	return doRequest[FactionCategory]("PUT", fmt.Sprintf("/api/faction-categories/%d", id), body)
}

func DeleteFactionCategory(id int) error {
	_, err := doRequest[any]("DELETE", fmt.Sprintf("/api/faction-categories/%d", id), nil)
	return err
}

func ExportWorld() ([]byte, error) {
	req, err := http.NewRequest("GET", BaseURL+"/admin/export", nil)
	if err != nil {
		return nil, err
	}
	if token := GetToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// ─── Admin ───────────────────────────────────────────────────────────────────

func WipeWorld() error {
	_, err := doRequest[any]("POST", "/admin/wipe/full", nil)
	return err
}

// must panics on error (used for constants)
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
