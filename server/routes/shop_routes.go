package routes

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/dblog"
	"herbst-server/middleware"
	"herbst-server/repository"
)

// RegisterShopRoutes registers REST endpoints for shop templates and items.
func RegisterShopRoutes(r *gin.Engine, repos *repository.Container) {
	shops := r.Group("/api")
	shops.Use(middleware.AuthMiddleware(nil))
	shops.Use(middleware.AdminMiddleware())
	shops.Use(middleware.WorldAccessMiddleware())
	{
		// Shop template CRUD
		shops.GET("/shops", listShopTemplates(repos))
		shops.POST("/shops", createShopTemplate(repos))
		shops.GET("/shops/:id", getShopTemplate(repos))
		shops.PUT("/shops/:id", updateShopTemplate(repos))
		shops.DELETE("/shops/:id", deleteShopTemplate(repos))

		// Shop items (nested under shop)
		shops.GET("/shops/:id/items", listShopItems(repos))
		shops.POST("/shops/:id/items", createShopItem(repos))
		shops.PUT("/shops/:id/items/:eqid", updateShopItem(repos))
		shops.DELETE("/shops/:id/items/:eqid", deleteShopItem(repos))
	}
}

// ─── Shop Template Views ─────────────────────────────────────────────────────

type shopTemplateView struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	WorldID           string `json:"world_id"`
	NPCTemplateID     string `json:"npc_template_id"`
	CurrencyItemType  int    `json:"currency_item_type"`
	MaxInventory      int    `json:"max_inventory"`
	GoldReserves      int    `json:"gold_reserves"`
	IsActive          bool   `json:"is_active"`
}

func toShopTemplateView(t *db.ShopTemplate) shopTemplateView {
	return shopTemplateView{
		ID:               t.ID,
		Name:             t.Name,
		WorldID:          t.WorldID,
		NPCTemplateID:    t.NpcTemplateID,
		CurrencyItemType: t.CurrencyItemType,
		MaxInventory:     t.MaxInventory,
		GoldReserves:     t.GoldReserves,
		IsActive:         t.IsActive,
	}
}

// ─── Shop Template Handlers ──────────────────────────────────────────────────

func listShopTemplates(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		worldID := c.Query("world_id")
		templates, err := repos.ShopTemplate.List(c.Request.Context(), worldID)
		if err != nil {
			dblog.Error("failed to list shop templates", err, slog.String("service", "shops"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]shopTemplateView, len(templates))
		for i, t := range templates {
			result[i] = toShopTemplateView(t)
		}
		c.JSON(http.StatusOK, result)
	}
}

func createShopTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name              string `json:"name"`
			WorldID           string `json:"world_id"`
			NPCTemplateID     string `json:"npc_template_id"`
			CurrencyItemType  int    `json:"currency_item_type"`
			MaxInventory      int    `json:"max_inventory"`
			GoldReserves      int    `json:"gold_reserves"`
			IsActive          bool   `json:"is_active"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if req.WorldID == "" {
			req.WorldID = c.Query("world_id")
		}
		if req.MaxInventory == 0 {
			req.MaxInventory = 50
		}
		if req.GoldReserves == 0 {
			req.GoldReserves = 1000
		}

		created, err := repos.ShopTemplate.Create(c.Request.Context(), repository.CreateShopTemplateInput{
			Name:             req.Name,
			WorldID:          req.WorldID,
			NPCTemplateID:    req.NPCTemplateID,
			CurrencyItemType: req.CurrencyItemType,
			MaxInventory:     req.MaxInventory,
			GoldReserves:     req.GoldReserves,
			IsActive:         req.IsActive,
		})
		if err != nil {
			dblog.Error("failed to create shop template", err, slog.String("service", "shops"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("shop template created", slog.Int("shop_id", created.ID), slog.String("service", "shops"))
		c.JSON(http.StatusCreated, toShopTemplateView(created))
	}
}

func getShopTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
			return
		}
		tmpl, err := repos.ShopTemplate.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "shop template not found"})
			return
		}
		c.JSON(http.StatusOK, toShopTemplateView(tmpl))
	}
}

func updateShopTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
			return
		}
		var req repository.ShopTemplateUpdates
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updated, err := repos.ShopTemplate.Update(c.Request.Context(), id, req)
		if err != nil {
			dblog.Error("failed to update shop template", err, slog.Int("shop_id", id))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("shop template updated", slog.Int("shop_id", updated.ID), slog.String("service", "shops"))
		c.JSON(http.StatusOK, toShopTemplateView(updated))
	}
}

func deleteShopTemplate(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
			return
		}
		if err := repos.ShopTemplate.Delete(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "shop template not found"})
			return
		}
		slog.Info("shop template deleted", slog.Int("shop_id", id), slog.String("service", "shops"))
		c.JSON(http.StatusNoContent, nil)
	}
}

// ─── Shop Item Views ─────────────────────────────────────────────────────────

type shopItemView struct {
	ID                  int    `json:"id"`
	ShopID              int    `json:"shop_id"`
	EquipmentTemplateID int    `json:"equipment_template_id"`
	Category            string `json:"category"`
	Price               int    `json:"price"`
	Quantity            int    `json:"quantity"`
	MaxStock            int    `json:"max_stock"`
	IsEnabled           bool   `json:"is_enabled"`
}

func toShopItemView(item *db.ShopItem) shopItemView {
	return shopItemView{
		ID:                  item.ID,
		ShopID:              item.ShopID,
		EquipmentTemplateID: item.EquipmentTemplateID,
		Category:            item.Category,
		Price:               item.Price,
		Quantity:            item.Quantity,
		MaxStock:            item.MaxStock,
		IsEnabled:           item.IsEnabled,
	}
}

// ─── Shop Item Handlers ──────────────────────────────────────────────────────

func listShopItems(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
			return
		}
		items, err := repos.ShopItem.ListByShop(c.Request.Context(), shopID)
		if err != nil {
			dblog.Error("failed to list shop items", err, slog.Int("shop_id", shopID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]shopItemView, len(items))
		for i, item := range items {
			result[i] = toShopItemView(item)
		}
		c.JSON(http.StatusOK, result)
	}
}

func createShopItem(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
			return
		}
		var req struct {
			EquipmentTemplateID int    `json:"equipment_template_id"`
			Category            string `json:"category"`
			Price               int    `json:"price"`
			Quantity            int    `json:"quantity"`
			MaxStock            int    `json:"max_stock"`
			IsEnabled           bool   `json:"is_enabled"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.EquipmentTemplateID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "equipment_template_id is required"})
			return
		}
		if req.MaxStock == 0 {
			req.MaxStock = 99
		}

		created, err := repos.ShopItem.Create(c.Request.Context(), repository.CreateShopItemInput{
			ShopID:              shopID,
			EquipmentTemplateID: req.EquipmentTemplateID,
			Category:            req.Category,
			Price:               req.Price,
			Quantity:            req.Quantity,
			MaxStock:            req.MaxStock,
			IsEnabled:           req.IsEnabled,
		})
		if err != nil {
			dblog.Error("failed to create shop item", err, slog.Int("shop_id", shopID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		slog.Info("shop item created", slog.Int("shop_id", shopID), slog.Int("equipment_template_id", req.EquipmentTemplateID), slog.String("service", "shops"))
		c.JSON(http.StatusCreated, toShopItemView(created))
	}
}

func updateShopItem(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
			return
		}
		eqID, err := strconv.Atoi(c.Param("eqid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid equipment template id"})
			return
		}
		var req repository.ShopItemUpdates
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updated, err := repos.ShopItem.Update(c.Request.Context(), shopID, eqID, req)
		if err != nil {
			dblog.Error("failed to update shop item", err, slog.Int("shop_id", shopID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, toShopItemView(updated))
	}
}

func deleteShopItem(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		shopID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop id"})
			return
		}
		eqID, err := strconv.Atoi(c.Param("eqid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid equipment template id"})
			return
		}
		if err := repos.ShopItem.Delete(c.Request.Context(), shopID, eqID); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "shop item not found"})
			return
		}
		slog.Info("shop item deleted", slog.Int("shop_id", shopID), slog.Int("equipment_template_id", eqID), slog.String("service", "shops"))
		c.JSON(http.StatusNoContent, nil)
	}
}

// ─── Public: Get shop by NPC template ID (for game client) ───────────────────

// GetShopByNPC is an unauthenticated endpoint used by the game client
// to fetch a shop's contents by NPC template ID.
func GetShopByNPC(repos *repository.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		npcTemplateID := c.Query("npc_template_id")
		if npcTemplateID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "npc_template_id is required"})
			return
		}

		ctx := context.Background()
		shop, err := repos.ShopTemplate.GetByNPCTemplateID(ctx, npcTemplateID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no shop found for this NPC"})
			return
		}

		items, err := repos.ShopItem.ListByShop(ctx, shop.ID)
		if err != nil {
			dblog.Error("failed to list shop items for client", err, slog.Int("shop_id", shop.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load shop items"})
			return
		}

		// Enrich items with equipment template details
		type clientShopItem struct {
			ID                  int    `json:"id"`
			ShopID              int    `json:"shop_id"`
			EquipmentTemplateID int    `json:"equipment_template_id"`
			Name                string `json:"name"`
			Description         string `json:"description"`
			Category            string `json:"category"`
			Price               int    `json:"price"`
			Quantity            int    `json:"quantity"`
			MaxStock            int    `json:"max_stock"`
			IsEnabled           bool   `json:"is_enabled"`
			Slot                string `json:"slot"`
			Rarity              string `json:"rarity"`
			DamageDiceCount     int    `json:"damage_dice_count"`
			DamageDiceSides     int    `json:"damage_dice_sides"`
			DamageBonus         int    `json:"damage_bonus"`
			DamageType          string `json:"damage_type"`
		}

		clientItems := make([]clientShopItem, 0, len(items))
		for _, item := range items {
			if !item.IsEnabled || item.Quantity <= 0 {
				continue
			}
			eq, err := repos.EquipmentTemplate.Get(ctx, item.EquipmentTemplateID)
			if err != nil {
				continue
			}
			clientItems = append(clientItems, clientShopItem{
				ID:                  item.ID,
				ShopID:              item.ShopID,
				EquipmentTemplateID: item.EquipmentTemplateID,
				Name:                eq.Name,
				Description:         eq.Description,
				Category:            item.Category,
				Price:               item.Price,
				Quantity:            item.Quantity,
				MaxStock:            item.MaxStock,
				IsEnabled:           item.IsEnabled,
				Slot:                eq.Slot,
				Rarity:              eq.Rarity,
				DamageDiceCount:     eq.DamageDiceCount,
				DamageDiceSides:     eq.DamageDiceSides,
				DamageBonus:         eq.DamageBonus,
				DamageType:          eq.DamageType,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"id":            shop.ID,
			"name":          shop.Name,
			"npc_template_id": shop.NpcTemplateID,
			"currency_name": "gold",
			"gold":          shop.GoldReserves,
			"items":         clientItems,
		})
	}
}