package routes

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/service"
)

// RegisterQuestRoutes registers CRUD endpoints for Quest definitions.
// All quest routes require admin authentication + world access.
func RegisterQuestRoutes(r *gin.Engine, svc *service.Container) {
	quests := r.Group("/api/quests")
	quests.Use(middleware.AuthMiddleware(nil))
	quests.Use(middleware.AdminMiddleware())
	quests.Use(middleware.WorldAccessMiddleware())
	{
		quests.GET("", listQuests(svc))
		quests.POST("", createQuest(svc))
		quests.GET("/:id", getQuest(svc))
		quests.PUT("/:id", updateQuest(svc))
		quests.DELETE("/:id", deleteQuest(svc))
		quests.GET("/lookups", getQuestLookups(svc))
	}
}

// questLookups holds all lookup data for quest form fields.
type questLookups struct {
	QuestTypes        []lookupItem        `json:"quest_types"`
	NPCs              []lookupItem        `json:"npcs"`
	Rooms             []lookupItem        `json:"rooms"`
	Items             []lookupItem        `json:"items"`
	Effects           []lookupItem        `json:"effects"`
	Tags              []lookupItem        `json:"tags"`
	Achievements      []lookupItem        `json:"achievements"`
	PrerequisiteQuest []lookupItem        `json:"prerequisite_quests"`
}

type lookupItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// getQuestLookups returns all lookup data needed for quest form fields.
func getQuestLookups(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Quest types: kill, explore, collect
		questTypes := []lookupItem{
			{ID: "kill", Name: "Kill"},
			{ID: "explore", Name: "Explore"},
			{ID: "collect", Name: "Collect"},
		}

		// NPCs for kill targets
		npcs, err := svc.NPC.ListTemplates(ctx, "")
		if err != nil {
			dblog.Error("failed to load NPCs", err, slog.String("service", "quests"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load NPCs"})
			return
		}
		npcItems := make([]lookupItem, len(npcs))
		for i, n := range npcs {
			npcItems[i] = lookupItem{ID: n.ID, Name: n.Name}
		}
		sort.Slice(npcItems, func(i, j int) bool { return npcItems[i].Name < npcItems[j].Name })

		// Rooms for explore targets
		rooms, err := svc.Room.ListRooms(ctx, "")
		if err != nil {
			dblog.Error("failed to load rooms", err, slog.String("service", "quests"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load rooms"})
			return
		}
		roomItems := make([]lookupItem, len(rooms))
		for i, r := range rooms {
			roomItems[i] = lookupItem{ID: fmt.Sprintf("%d", r.ID), Name: r.Name}
		}
		sort.Slice(roomItems, func(i, j int) bool { return roomItems[i].Name < roomItems[j].Name })

		// Effects - use EffectRepo directly
		effects, err := svc.Client.Effect.Query().All(ctx)
		if err != nil {
			dblog.Error("failed to load effects", err, slog.String("service", "quests"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load effects"})
			return
		}
		effectItems := make([]lookupItem, len(effects))
		for i, e := range effects {
			effectItems[i] = lookupItem{ID: fmt.Sprintf("%d", e.ID), Name: e.Name}
		}
		sort.Slice(effectItems, func(i, j int) bool { return effectItems[i].Name < effectItems[j].Name })

		// Tags - use TagRepo directly
		tags, err := svc.Client.Tag.Query().All(ctx)
		if err != nil {
			dblog.Error("failed to load tags", err, slog.String("service", "quests"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tags"})
			return
		}
		tagItems := make([]lookupItem, len(tags))
		for i, t := range tags {
			tagItems[i] = lookupItem{ID: fmt.Sprintf("%d", t.ID), Name: t.Name}
		}
		sort.Slice(tagItems, func(i, j int) bool { return tagItems[i].Name < tagItems[j].Name })

		// Achievements - use AchievementRepo directly
		achievements, err := svc.Client.Achievement.Query().All(ctx)
		if err != nil {
			dblog.Error("failed to load achievements", err, slog.String("service", "quests"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load achievements"})
			return
		}
		achievementItems := make([]lookupItem, len(achievements))
		for i, a := range achievements {
			achievementItems[i] = lookupItem{ID: fmt.Sprintf("%d", a.ID), Name: a.Name}
		}
		sort.Slice(achievementItems, func(i, j int) bool { return achievementItems[i].Name < achievementItems[j].Name })

		// Items - from EquipmentTemplate (item definitions)
		items, err := svc.Client.EquipmentTemplate.Query().All(ctx)
		if err != nil {
			dblog.Error("failed to load items", err, slog.String("service", "quests"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load items"})
			return
		}
		itemItems := make([]lookupItem, len(items))
		for i, item := range items {
			itemItems[i] = lookupItem{ID: fmt.Sprintf("%d", item.ID), Name: item.Name}
		}
		sort.Slice(itemItems, func(i, j int) bool { return itemItems[i].Name < itemItems[j].Name })

		// Prerequisite quests
		allQuests, err := svc.Quest.ListQuests(ctx, "")
		if err != nil {
			dblog.Error("failed to load quests", err, slog.String("service", "quests"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load quests"})
			return
		}
		prereqQuestItems := make([]lookupItem, len(allQuests))
		for i, q := range allQuests {
			prereqQuestItems[i] = lookupItem{ID: fmt.Sprintf("%d", q.ID), Name: q.Name}
		}
		sort.Slice(prereqQuestItems, func(i, j int) bool { return prereqQuestItems[i].Name < prereqQuestItems[j].Name })

		c.JSON(http.StatusOK, questLookups{
			QuestTypes:         questTypes,
			NPCs:               npcItems,
			Rooms:              roomItems,
			Items:              itemItems,
			Effects:            effectItems,
			Tags:               tagItems,
			Achievements:       achievementItems,
			PrerequisiteQuest:  prereqQuestItems,
		})
	}
}
