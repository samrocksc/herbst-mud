package dbinit

import (
	"context"
	"log"

	"herbst-server/db"
	"herbst-server/db/schema"
)

// InitCraftingContent seeds Pizza Chef faction, equipment templates, and recipes.
func InitCraftingContent(client *db.Client) error {
	ctx := context.Background()

	// Seed faction category "class"
	cats, _ := client.FactionCategory.Query().Where().All(ctx)
	hasClass := false
	for _, c := range cats {
		if c.Name == "class" {
			hasClass = true
			break
		}
	}
	var classCatID int
	if !hasClass {
		cat, err := client.FactionCategory.Create().
			SetName("class").
			SetDisplayName("Class").
			SetDescription("Character class factions").
			SetMaxMemberships(1).
			SetInitialConfig(true).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed faction category 'class': %v", err)
		} else {
			classCatID = cat.ID
			log.Println("Faction category 'class' seeded")
		}
	} else {
		for _, c := range cats {
			if c.Name == "class" {
				classCatID = c.ID
				break
			}
		}
	}

	// Seed pizza_chef faction
	factions, _ := client.Faction.Query().Where().All(ctx)
	hasPizzaChef := false
	for _, f := range factions {
		if f.Name == "pizza_chef" {
			hasPizzaChef = true
			break
		}
	}
	if !hasPizzaChef {
		faction, err := client.Faction.Create().
			SetName("pizza_chef").
			SetWorldID("default").
			SetDisplayName("Pizza Chef").
			SetDescription("Those who craft delicious pizzas.").
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed faction 'pizza_chef': %v", err)
		} else {
			// Assign category via edge mutation
			if classCatID > 0 {
				cat, _ := client.FactionCategory.Get(ctx, classCatID)
				if cat != nil {
					_, err = client.Faction.UpdateOne(faction).SetCategory(cat).Save(ctx)
					if err != nil {
						log.Printf("Warning: failed to assign faction category: %v", err)
					}
				}
			}
			log.Println("Faction 'pizza_chef' seeded")
		}
	}

	// Seed ingredient equipment templates (misc type)
	ingredients := []struct {
		id          string
		name        string
		description string
		itemType    string
		weight      int
	}{
		{"dough", "dough", "Fresh pizza dough", "misc", 1},
		{"sauce", "sauce", "Tangy tomato sauce", "misc", 1},
		{"cheese", "cheese", "Mozzarella cheese", "misc", 1},
		{"pepperoni", "pepperoni", "Spicy pepperoni slices", "misc", 1},
		{"mushroom", "mushroom", "Sliced mushrooms", "misc", 1},
		{"olive", "olive", "Black olives", "misc", 1},
	}

	existingTemplates, _ := client.EquipmentTemplate.Query().Where().All(ctx)
	existingIDs := make(map[string]bool)
	for _, t := range existingTemplates {
		existingIDs[t.ID] = true
	}

	for _, ing := range ingredients {
		if existingIDs[ing.id] {
			continue
		}
		_, err := client.EquipmentTemplate.Create().
			SetID(ing.id).
			SetWorldID("default").
			SetName(ing.name).
			SetDescription(ing.description).
			SetSlot("misc").
			SetItemType(ing.itemType).
			SetWeight(ing.weight).
			SetIsVisible(true).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed equipment template '%s': %v", ing.id, err)
		} else {
			log.Printf("Equipment template '%s' seeded", ing.id)
		}
	}

	// Seed pizza weapon templates
	pizzas := []struct {
		id              string
		name            string
		description     string
		damageDiceCount int
		damageDiceSides int
		damageBonus     int
		damageType      string
		weaponType      string
		isTwoHanded     bool
		weight          int
	}{
		{"cheese_pizza", "cheese_pizza", "Classic cheese pizza.", 1, 3, 1, "bludgeoning", "improvised", false, 2},
		{"pepperoni_pizza", "pepperoni_pizza", "A delicious pepperoni pizza. Still warm.", 1, 4, 2, "bludgeoning", "improvised", false, 2},
		{"supreme_pizza", "supreme_pizza", "Loaded supreme pizza with everything on it.", 1, 6, 3, "bludgeoning", "improvised", false, 3},
	}

	for _, p := range pizzas {
		if existingIDs[p.id] {
			continue
		}
		_, err := client.EquipmentTemplate.Create().
			SetID(p.id).
			SetWorldID("default").
			SetName(p.name).
			SetDescription(p.description).
			SetSlot("main_hand").
			SetItemType("weapon").
			SetWeight(p.weight).
			SetDamageDiceCount(p.damageDiceCount).
			SetDamageDiceSides(p.damageDiceSides).
			SetDamageBonus(p.damageBonus).
			SetDamageType(p.damageType).
			SetWeaponType(p.weaponType).
			SetIsTwoHanded(p.isTwoHanded).
			SetIsVisible(true).
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed weapon template '%s': %v", p.id, err)
		} else {
			log.Printf("Weapon template '%s' seeded", p.id)
		}
	}

	// Seed crafting recipes
	recipes := []struct {
		name          string
		displayName   string
		description   string
		stationTag    string
		class         string
		inputs        []schema.CraftingInput
		outputs       []schema.CraftingOutput
		craftTimeSecs int
	}{
		{
			name:        "cheese_pizza_recipe",
			displayName: "Cheese Pizza",
			description: "Classic cheese pizza.",
			stationTag:  "pizza_station",
			class:       "pizza_chef",
			inputs: []schema.CraftingInput{
				{EquipmentTemplateID: "dough", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "sauce", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "cheese", Quantity: 1, Consumed: true},
			},
			outputs: []schema.CraftingOutput{
				{EquipmentTemplateID: "cheese_pizza", Quantity: 1},
			},
			craftTimeSecs: 3,
		},
		{
			name:        "pepperoni_pizza_recipe",
			displayName: "Pepperoni Pizza",
			description: "Pepperoni pizza with spicy slices.",
			stationTag:  "pizza_station",
			class:       "pizza_chef",
			inputs: []schema.CraftingInput{
				{EquipmentTemplateID: "dough", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "sauce", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "cheese", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "pepperoni", Quantity: 1, Consumed: true},
			},
			outputs: []schema.CraftingOutput{
				{EquipmentTemplateID: "pepperoni_pizza", Quantity: 1},
			},
			craftTimeSecs: 3,
		},
		{
			name:        "supreme_pizza_recipe",
			displayName: "Supreme Pizza",
			description: "Loaded supreme with all toppings.",
			stationTag:  "pizza_station",
			class:       "pizza_chef",
			inputs: []schema.CraftingInput{
				{EquipmentTemplateID: "dough", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "sauce", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "cheese", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "pepperoni", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "mushroom", Quantity: 1, Consumed: true},
				{EquipmentTemplateID: "olive", Quantity: 1, Consumed: true},
			},
			outputs: []schema.CraftingOutput{
				{EquipmentTemplateID: "supreme_pizza", Quantity: 1},
			},
			craftTimeSecs: 5,
		},
	}

	existingRecipes, _ := client.CraftingRecipe.Query().Where().All(ctx)
	existingRecipeNames := make(map[string]bool)
	for _, r := range existingRecipes {
		existingRecipeNames[r.Name] = true
	}

	for _, r := range recipes {
		if existingRecipeNames[r.name] {
			continue
		}
		_, err := client.CraftingRecipe.Create().
			SetName(r.name).
			SetDisplayName(r.displayName).
			SetDescription(r.description).
			SetRequiredStationTag(r.stationTag).
			SetRequiredClass(r.class).
			SetInputs(r.inputs).
			SetOutputs(r.outputs).
			SetCraftTimeSecs(r.craftTimeSecs).
			SetWorldID("default").
			Save(ctx)
		if err != nil {
			log.Printf("Warning: failed to seed recipe '%s': %v", r.name, err)
		} else {
			log.Printf("Recipe '%s' seeded", r.name)
		}
	}

	// Tag room 1 with pizza_station
	room, err := client.Room.Get(ctx, 1)
	if err == nil && room != nil {
		currentTags := room.Tags
		hasStation := false
		for _, t := range currentTags {
			if t == "pizza_station" {
				hasStation = true
				break
			}
		}
		if !hasStation {
			newTags := append(currentTags, "pizza_station")
			_, err = client.Room.UpdateOne(room).SetTags(newTags).Save(ctx)
			if err != nil {
				log.Printf("Warning: failed to tag room 1 with pizza_station: %v", err)
			} else {
				log.Println("Room 1 tagged with pizza_station")
			}
		} else {
			log.Println("Room 1 already has pizza_station tag")
		}
	} else {
		log.Printf("Warning: room 1 not found: %v", err)
	}

	log.Println("Crafting content seeded successfully")
	return nil
}