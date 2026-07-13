package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"

	"entgo.io/ent/dialect/sql"
	"herbst-server/db"
	"herbst-server/db/character"
	"herbst-server/db/charactercompetency"
	"herbst-server/db/characterfaction"
	"herbst-server/db/competencycategory"
	"herbst-server/db/competencylevelthreshold"
	"herbst-server/db/gameconfig"
	"herbst-server/db/race"
	"herbst-server/db/world"
	"herbst-server/events"
)

// DefaultXPPerLevel is used when no GameConfig override exists.
const DefaultXPPerLevel = 200

// Default stat growth values (per level) used when world config is absent.
const (
	DefaultHPPerLevel     = 10
	DefaultManaPerLevel   = 5
	DefaultStaminaPerLevel = 5
)

// Default anti-grind kill threshold.
const DefaultAntiGrindKillThreshold = 20

// XPThresholds holds the XP required to reach each level.
type XPThresholds map[int]int

// XPThresholdsFromConfig parses thresholds from a GameConfig value JSON string.
func XPThresholdsFromConfig(value string) (XPThresholds, error) {
	if value == "" {
		return nil, fmt.Errorf("empty config value")
	}
	var t XPThresholds
	if err := json.Unmarshal([]byte(value), &t); err != nil {
		return nil, err
	}
	return t, nil
}

// ThresholdForLevel returns the XP required to reach the given level.
// Returns 0 for level < 1.
func (t XPThresholds) ThresholdForLevel(level int) int {
	if level < 1 {
		return 0
	}
	if xp, ok := t[level]; ok {
		return xp
	}
	// Fallback: DefaultXPPerLevel * level
	return DefaultXPPerLevel * level
}

// xpAwardService handles XP awards and level-up logic.
type xpAwardService struct {
	client *db.Client
	logger *slog.Logger
}

// NewXPAwardService creates a new XP service.
func NewXPAwardService(client *db.Client, logger *slog.Logger) XPAwardService {
	if logger == nil {
		logger = slog.Default()
	}
	return &xpAwardService{client: client, logger: logger}
}

// AwardXP adds XP to a character and handles level-up if thresholds are crossed.
// Returns new XP, new level, and whether a level-up occurred.
func (s *xpAwardService) AwardXP(ctx context.Context, characterID, xpGained int) (newXP, newLevel int, leveledUp bool, err error) {
	return s.AwardXPWithSource(ctx, characterID, xpGained, "kill")
}

// AwardXPWithSource is like AwardXP but lets the caller specify the XP source
// (e.g. "kill", "quest", "exploration") for the emitted xp.gained event.
func (s *xpAwardService) AwardXPWithSource(ctx context.Context, characterID, xpGained int, source string) (newXP, newLevel int, leveledUp bool, err error) {
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, false, fmt.Errorf("get character: %w", err)
	}

	charXP := char.Xp
	charLevel := char.Level

	// Apply XP
	charXP += xpGained

	// Load world config for level curve
	worldConfig, worldCfgMap := s.loadWorldConfig(ctx, char)

	// Get thresholds (world config level curve or fallback)
	thresholds, maxLevel := s.getThresholdsForCharacter(ctx, char, worldCfgMap)

	// Check for level-up (can happen multiple times in one award)
	for {
		if maxLevel > 0 && charLevel >= maxLevel {
			break
		}
		nextLevel := charLevel + 1
		needed := thresholds.ThresholdForLevel(nextLevel)
		if needed == 0 || charXP < needed {
			break
		}
		charLevel = nextLevel
		s.logger.Info("level up",
			"character_id", characterID,
			"new_level", charLevel,
			"xp_remaining", charXP,
		)
	}

	// Update character XP and level
	upd := s.client.Character.UpdateOne(char).
		SetXp(charXP).
		SetLevel(charLevel)

	// Scale stats on level-up
	if charLevel > char.Level {
		hpGain, manaGain, staminaGain := s.computeStatGrowth(ctx, char, charLevel-char.Level, worldConfig, worldCfgMap)
		if hpGain > 0 || manaGain > 0 || staminaGain > 0 {
			newMaxHP := char.MaxHitpoints + hpGain
			newMaxMana := char.MaxMana + manaGain
			newMaxStamina := char.MaxStamina + staminaGain
			upd = upd.
				SetMaxHitpoints(newMaxHP).
				SetMaxMana(newMaxMana).
				SetMaxStamina(newMaxStamina).
				SetHitpoints(newMaxHP).
				SetMana(newMaxMana).
				SetStamina(newMaxStamina)
			s.logger.Info("stat growth on level up",
				"character_id", characterID,
				"hp_gain", hpGain,
				"mana_gain", manaGain,
				"stamina_gain", staminaGain,
				"new_max_hp", newMaxHP,
				"new_max_mana", newMaxMana,
				"new_max_stamina", newMaxStamina,
			)
		}
	}

	_, err = upd.Save(ctx)
	if err != nil {
		return 0, 0, false, fmt.Errorf("update character xp/level: %w", err)
	}

	leveledUp = charLevel > char.Level

	// Emit xp.gained event
	events.Publish(events.Event{
		Type: events.EventXPGained,
		Payload: map[string]interface{}{
			"character_id": characterID,
			"amount":        xpGained,
			"new_xp":        charXP,
			"new_level":     charLevel,
			"source":        source,
		},
	})

	return charXP, charLevel, leveledUp, nil
}

// loadWorldConfig loads the world config map from the character's world.
// Returns the parsed config (as a typed struct) and the raw map for misc lookups.
func (s *xpAwardService) loadWorldConfig(ctx context.Context, char *db.Character) (WorldConfig, map[string]interface{}) {
	var wc WorldConfig
	var raw map[string]interface{}

	if char.WorldID == 0 {
		return wc, nil
	}

	w, err := s.client.World.Query().Where(world.IDEQ(char.WorldID)).Only(ctx)
	if err != nil {
		s.logger.Debug("failed to load world for config", "world_id", char.WorldID, "error", err)
		return wc, nil
	}

	if w.Config == nil || len(w.Config) == 0 {
		return wc, nil
	}

	raw = w.Config

	// Parse level_curve
	if lc, ok := raw["level_curve"].(map[string]interface{}); ok {
		if mode, ok := lc["mode"].(string); ok {
			wc.LevelCurve.Mode = mode
		}
		if baseXP, ok := lc["base_xp"].(float64); ok {
			wc.LevelCurve.BaseXP = int(baseXP)
		}
		if pct, ok := lc["percentage"].(float64); ok {
			wc.LevelCurve.Percentage = pct
		}
		if maxLvl, ok := lc["max_level"].(float64); ok {
			wc.LevelCurve.MaxLevel = int(maxLvl)
		}
		if thresholds, ok := lc["thresholds"].([]interface{}); ok {
			for _, t := range thresholds {
				if tf, ok := t.(float64); ok {
					wc.LevelCurve.Thresholds = append(wc.LevelCurve.Thresholds, int(tf))
				}
			}
		}
	}

	// Parse stat_growth
	if sg, ok := raw["stat_growth"].(map[string]interface{}); ok {
		if hp, ok := sg["hp_per_level"].(float64); ok {
			wc.StatGrowth.HPPerLevel = int(hp)
		}
		if mana, ok := sg["mana_per_level"].(float64); ok {
			wc.StatGrowth.ManaPerLevel = int(mana)
		}
		if stam, ok := sg["stamina_per_level"].(float64); ok {
			wc.StatGrowth.StaminaPerLevel = int(stam)
		}
	}

	// Parse skill_xp
	if sx, ok := raw["skill_xp"].(map[string]interface{}); ok {
		if t, ok := sx["anti_grind_kill_threshold"].(float64); ok {
			wc.SkillXP.AntiGrindKillThreshold = int(t)
		}
	}

	return wc, raw
}

// WorldConfig is a typed view of the world config JSON.
type WorldConfig struct {
	LevelCurve struct {
		Mode       string  `json:"mode"`
		BaseXP     int     `json:"base_xp"`
		Percentage float64 `json:"percentage"`
		MaxLevel   int     `json:"max_level"`
		Thresholds []int   `json:"thresholds"`
	} `json:"level_curve"`
	StatGrowth struct {
		HPPerLevel     int `json:"hp_per_level"`
		ManaPerLevel   int `json:"mana_per_level"`
		StaminaPerLevel int `json:"stamina_per_level"`
	} `json:"stat_growth"`
	SkillXP struct {
		AntiGrindKillThreshold int `json:"anti_grind_kill_threshold"`
	} `json:"skill_xp"`
}

// getThresholdsForCharacter returns XP thresholds and max_level for the character.
// Uses world config level curve if available, otherwise falls back to GameConfig.
func (s *xpAwardService) getThresholdsForCharacter(ctx context.Context, char *db.Character, worldCfg map[string]interface{}) (XPThresholds, int) {
	if worldCfg != nil {
		if lc, ok := worldCfg["level_curve"].(map[string]interface{}); ok {
			mode, _ := lc["mode"].(string)
			switch mode {
			case "percentage":
				baseXP, _ := lc["base_xp"].(float64)
				if baseXP == 0 {
					baseXP = float64(DefaultXPPerLevel)
				}
				pct, _ := lc["percentage"].(float64)
				if pct == 0 {
					pct = 50
				}
				maxLvl, _ := lc["max_level"].(float64)
				maxLevel := int(maxLvl)
				if maxLevel == 0 {
					maxLevel = 50
				}
				thresholds := make(XPThresholds)
				cumulative := 0
				for i := 1; i <= maxLevel; i++ {
					cumulative += int(math.Round(baseXP * math.Pow(1+pct/100, float64(i-1))))
					thresholds[i+1] = cumulative
				}
				return thresholds, maxLevel
			case "hand_coded":
				thresholdsArr, _ := lc["thresholds"].([]interface{})
				thresholds := make(XPThresholds)
				for i, t := range thresholdsArr {
					if tf, ok := t.(float64); ok {
						thresholds[i+2] = int(tf)
					}
				}
				maxLvl, _ := lc["max_level"].(float64)
				return thresholds, int(maxLvl)
			}
		}
	}

	// Fallback: use existing GameConfig-based approach
	return s.getThresholds(ctx), 0
}

// computeStatGrowth computes HP/mana/stamina gain for the given number of levels.
// Formula per level: base × class_modifier × racial_multiplier
func (s *xpAwardService) computeStatGrowth(ctx context.Context, char *db.Character, levelsGained int, wc WorldConfig, worldCfg map[string]interface{}) (hpGain, manaGain, staminaGain int) {
	if levelsGained <= 0 {
		return 0, 0, 0
	}

	// Base growth from world config (defaults if absent)
	baseHP := wc.StatGrowth.HPPerLevel
	if baseHP == 0 {
		baseHP = DefaultHPPerLevel
	}
	baseMana := wc.StatGrowth.ManaPerLevel
	if baseMana == 0 {
		baseMana = DefaultManaPerLevel
	}
	baseStamina := wc.StatGrowth.StaminaPerLevel
	if baseStamina == 0 {
		baseStamina = DefaultStaminaPerLevel
	}

	// Class/faction modifier from faction.stat_bonuses
	// We use constitution bonus as HP modifier, intelligence/wisdom for mana, dexterity for stamina.
	// If no class faction is found, modifiers default to 1.0.
	hpClassMod, manaClassMod, staminaClassMod := s.getClassModifiers(ctx, char)

	// Racial multiplier from race.stat_growth_multipliers
	hpRaceMod, manaRaceMod, staminaRaceMod := s.getRaceMultipliers(ctx, char)

	for i := 0; i < levelsGained; i++ {
		hpGain += int(math.Round(float64(baseHP) * hpClassMod * hpRaceMod))
		manaGain += int(math.Round(float64(baseMana) * manaClassMod * manaRaceMod))
		staminaGain += int(math.Round(float64(baseStamina) * staminaClassMod * staminaRaceMod))
	}

	return hpGain, manaGain, staminaGain
}

// getClassModifiers queries the character's active class faction and derives
// HP/mana/stamina modifiers from its stat_bonuses.
// Constitution → HP modifier (1 + constitution/20), Intelligence → mana (1 + int/20), Dexterity → stamina (1 + dex/20)
func (s *xpAwardService) getClassModifiers(ctx context.Context, char *db.Character) (hpMod, manaMod, staminaMod float64) {
	hpMod = 1.0
	manaMod = 1.0
	staminaMod = 1.0

	// Query active character faction memberships with their faction edge loaded
	cfs, err := s.client.CharacterFaction.Query().
		Where(characterfaction.HasCharacterWith(character.ID(char.ID))).
		Where(characterfaction.StatusEQ("active")).
		WithFaction().
		All(ctx)
	if err != nil {
		s.logger.Debug("failed to query character factions for class modifier", "error", err)
		return
	}

	for _, cf := range cfs {
		f := cf.Edges.Faction
		if f == nil {
			continue
		}
		// Use stat_bonuses from this faction
		sb := f.StatBonuses
		if sb.Constitution > 0 {
			hpMod += float64(sb.Constitution) / 20.0
		}
		if sb.Intelligence > 0 {
			manaMod += float64(sb.Intelligence) / 20.0
		}
		if sb.Dexterity > 0 {
			staminaMod += float64(sb.Dexterity) / 20.0
		}
	}
	return
}

// getRaceMultipliers queries the character's race and returns the stat growth multipliers.
// Defaults to 1.0 for each if not specified.
func (s *xpAwardService) getRaceMultipliers(ctx context.Context, char *db.Character) (hpMult, manaMult, staminaMult float64) {
	hpMult = 1.0
	manaMult = 1.0
	staminaMult = 1.0

	if char.Race == "" {
		return
	}

	// Query race by name + world (CurrentWorld is the string world ID)
	r, err := s.client.Race.Query().
		Where(race.NameEQ(char.Race)).
		All(ctx)
	if err != nil || len(r) == 0 {
		// Try without world filter if not found
		s.logger.Debug("failed to query race for stat growth multipliers", "race", char.Race, "error", err)
		return
	}

	// Use first match (races are unique by name+world, but we queried without world filter)
	raceObj := r[0]
	if raceObj.StatGrowthMultipliers == nil {
		return
	}

	if v, ok := raceObj.StatGrowthMultipliers["hp"]; ok && v > 0 {
		hpMult = v
	}
	if v, ok := raceObj.StatGrowthMultipliers["mana"]; ok && v > 0 {
		manaMult = v
	}
	if v, ok := raceObj.StatGrowthMultipliers["stamina"]; ok && v > 0 {
		staminaMult = v
	}
	return
}

// ApplyDeathPenalty reduces XP by penaltyPercent and returns XP lost and new XP.
func (s *xpAwardService) ApplyDeathPenalty(ctx context.Context, characterID, penaltyPercent int) (xpLost, newXP int, err error) {
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, fmt.Errorf("get character: %w", err)
	}

	xpLost = (char.Xp * penaltyPercent) / 100
	newXP = char.Xp - xpLost

	_, err = s.client.Character.UpdateOne(char).
		SetXp(newXP).
		Save(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("update character xp after death: %w", err)
	}

	return xpLost, newXP, nil
}

// getThresholds loads XP thresholds from GameConfig or returns the default.
func (s *xpAwardService) getThresholds(ctx context.Context) XPThresholds {
	cfg, err := s.client.GameConfig.Query().Where(gameconfig.Key("xp_thresholds")).Only(ctx)
	if err != nil {
		// Fall back to default linear thresholds
		thresholds := make(XPThresholds)
		for i := 1; i <= 100; i++ {
			thresholds[i] = DefaultXPPerLevel * i
		}
		return thresholds
	}

	thresholds, err := XPThresholdsFromConfig(cfg.Value)
	if err != nil {
		s.logger.Warn("failed to parse xp_thresholds config, using defaults", "error", err)
		thresholds = make(XPThresholds)
		for i := 1; i <= 100; i++ {
			thresholds[i] = DefaultXPPerLevel * i
		}
	}
	return thresholds
}

// GetCharacterXP returns the current XP and level for a character.
func (s *xpAwardService) GetCharacterXP(ctx context.Context, characterID int) (xp, level int, err error) {
	char, err := s.client.Character.Get(ctx, characterID)
	if err != nil {
		return 0, 0, err
	}
	return char.Xp, char.Level, nil
}

// QueryCharacter queries a character by ID (used by other services).
func (s *xpAwardService) QueryCharacter(ctx context.Context, id int) (*db.Character, error) {
	return s.client.Character.Get(ctx, id)
}

// AwardCompetencyXP awards XP to a character's competency in a category.
// It applies the category's xp_multiplier, updates the character's competency record,
// and recomputes the cached level based on thresholds.
func (s *xpAwardService) AwardCompetencyXP(ctx context.Context, characterID int, categoryID string, rawXP int) error {
	cat, err := s.client.CompetencyCategory.Get(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("get competency category %s: %w", categoryID, err)
	}

	multiplied := int(float64(rawXP) * cat.XpMultiplier)

	cc, err := s.client.CharacterCompetency.Query().
		Where(charactercompetency.HasCharacterWith(character.ID(characterID))).
		Where(charactercompetency.HasCategoryWith(competencycategory.ID(categoryID))).
		Only(ctx)
	if err != nil {
		// Record doesn't exist — create it
		cc, err = s.client.CharacterCompetency.Create().
			SetXp(multiplied).
			SetLevel(1).
			SetCharacterID(characterID).
			SetCategoryID(categoryID).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create character competency: %w", err)
		}
		s.logger.Info("competency started",
			"character_id", characterID, "category", categoryID, "xp", multiplied)
		return nil
	}

	// Update XP
	cc.Xp += multiplied

	// Recompute level from thresholds
	thresholds, err := s.client.CompetencyLevelThreshold.Query().
		Where(competencylevelthreshold.HasCategoryWith(competencycategory.ID(categoryID))).
		Order(competencylevelthreshold.ByLevel(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query thresholds: %w", err)
	}

	newLevel := cc.Level
	for _, t := range thresholds {
		if cc.Xp >= t.XpRequired {
			newLevel = t.Level
		}
	}

	cc.Level = newLevel

	_, err = s.client.CharacterCompetency.UpdateOne(cc).SetXp(cc.Xp).SetLevel(cc.Level).Save(ctx)
	if err != nil {
		return fmt.Errorf("update character competency: %w", err)
	}

	s.logger.Info("competency xp awarded",
		"character_id", characterID, "category", categoryID,
		"raw_xp", rawXP, "multiplied", multiplied, "total_xp", cc.Xp, "level", cc.Level)
	return nil
}

// SeedCompetencyCategories creates the default competency categories and level thresholds
// if they don't already exist. Safe to call on every startup.
func (s *xpAwardService) SeedCompetencyCategories(ctx context.Context) error {
	categories := []struct {
		id   string
		name string
		mult float64
	}{
		{"blades", "Blades", 0.20},
		{"staves", "Staves", 0.20},
		{"knives", "Knives", 0.20},
		{"martial", "Martial Arts", 0.20},
		{"brawling", "Brawling", 0.20},
		{"tech", "Tech", 0.20},
		{"light_armor", "Light Armor", 0.20},
		{"cloth_armor", "Cloth Armor", 0.20},
		{"heavy_armor", "Heavy Armor", 0.20},
	}

	// Level thresholds: levels 1-10 with escalating XP costs
	levelXP := []int{0, 100, 250, 500, 900, 1500, 2300, 3300, 4500, 6000}
	dmgMults := []float64{1.0, 1.05, 1.10, 1.15, 1.20, 1.30, 1.40, 1.55, 1.70, 1.90}
	defMults := []float64{1.0, 1.03, 1.06, 1.09, 1.12, 1.16, 1.22, 1.30, 1.40, 1.55}

	for _, cat := range categories {
		existing, _ := s.client.CompetencyCategory.Get(ctx, cat.id)
		if existing != nil {
			continue
		}

		created, err := s.client.CompetencyCategory.Create().
			SetID(cat.id).
			SetName(cat.name).
			SetXpMultiplier(cat.mult).
			Save(ctx)
		if err != nil {
			s.logger.Warn("failed to seed competency category", "id", cat.id, "error", err)
			continue
		}

		// Create thresholds for levels 1-10
		for lvl := 1; lvl <= 10; lvl++ {
			_, err = s.client.CompetencyLevelThreshold.Create().
				SetID(fmt.Sprintf("%s-%d", cat.id, lvl)).
				SetLevel(lvl).
				SetXpRequired(levelXP[lvl-1]).
				SetDamageMultiplier(dmgMults[lvl-1]).
				SetDefenseMultiplier(defMults[lvl-1]).
				SetCategory(created).
				Save(ctx)
			if err != nil {
				s.logger.Warn("failed to seed threshold", "id", cat.id, "level", lvl, "error", err)
			}
		}
		s.logger.Info("seeded competency category", "id", cat.id, "name", cat.name)
	}
	return nil
}