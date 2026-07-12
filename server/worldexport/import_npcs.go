package worldexport

import (
	"context"
	"fmt"
	"strconv"

	"herbst-server/db"
)

func importNPCs(ctx context.Context, client *db.Client, npcs []map[string]interface{}, newWorldID string, maps *idMaps) (int, error) {
	worldIDInt, _ := strconv.Atoi(newWorldID)
	count := 0
	for _, n := range npcs {
		if !boolVal(n["isNPC"]) {
			continue
		}
		oldID := intVal(n["id"])
		currentRoom := maps.Rooms[intVal(n["currentRoomId"])]
		startingRoom := maps.Rooms[intVal(n["startingRoomId"])]
		respawnRoom := maps.Rooms[intVal(n["respawnRoomId"])]
		if currentRoom == 0 {
			currentRoom = firstRoom(maps)
		}
		if startingRoom == 0 {
			startingRoom = currentRoom
		}
		if respawnRoom == 0 {
			respawnRoom = currentRoom
		}

		created, err := client.Character.Create().
			SetName(strVal(n["name"])).
			SetIsNPC(true).
			SetCurrentRoomId(currentRoom).
			SetStartingRoomId(startingRoom).
			SetRespawnRoomId(respawnRoom).
			SetNillableWorldID(intPtr(worldIDInt)).
			SetCurrentWorld(newWorldID).
			SetRace(strValOr(n["race"], "human")).
			SetNillableClass(strPtr(n["class"])).
			SetLevel(intValOr(n["level"], 1)).
			SetHitpoints(intValOr(n["hitpoints"], 100)).
			SetMaxHitpoints(intValOr(n["max_hitpoints"], 100)).
			SetStamina(intValOr(n["stamina"], 50)).
			SetMaxStamina(intValOr(n["max_stamina"], 50)).
			SetMana(intValOr(n["mana"], 25)).
			SetMaxMana(intValOr(n["max_mana"], 25)).
			SetConstitution(intValOr(n["constitution"], 10)).
			SetStrength(intValOr(n["strength"], 10)).
			SetDexterity(intValOr(n["dexterity"], 10)).
			SetIntelligence(intValOr(n["intelligence"], 10)).
			SetWisdom(intValOr(n["wisdom"], 10)).
			SetNillableGender(strPtr(n["gender"])).
			SetNillableDescription(strPtr(n["description"])).
			SetNillableNpcTemplateID(strPtr(n["npc_template_id"])).
			SetNillableNpcSkillID(strPtr(n["npc_skill_id"])).
			SetXp(intValOr(n["xp"], 0)).
			SetGoldCredits(intValOr(n["gold_credits"], 0)).
			SetIsImmortal(boolVal(n["is_immortal"])).
			SetSkillBlades(intValOr(n["skill_blades"], 0)).
			SetSkillStaves(intValOr(n["skill_staves"], 0)).
			SetSkillKnives(intValOr(n["skill_knives"], 0)).
			SetSkillMartial(intValOr(n["skill_martial"], 0)).
			SetSkillBrawling(intValOr(n["skill_brawling"], 0)).
			SetSkillTech(intValOr(n["skill_tech"], 0)).
			SetSkillLightArmor(intValOr(n["skill_light_armor"], 0)).
			SetSkillClothArmor(intValOr(n["skill_cloth_armor"], 0)).
			SetSkillHeavyArmor(intValOr(n["skill_heavy_armor"], 0)).
			Save(ctx)
		if err != nil {
			return count, fmt.Errorf("npc %d: %w", oldID, err)
		}
		maps.Characters[oldID] = created.ID
		count++
	}
	return count, nil
}

func firstRoom(maps *idMaps) int {
	for _, v := range maps.Rooms {
		return v
	}
	return 0
}
