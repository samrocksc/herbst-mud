package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/db"
	"herbst-server/db/quest"
)

// listQuests returns all quests ordered by name.
func listQuests(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		quests, err := client.Quest.Query().
			Order(db.Asc(quest.FieldName)).
			All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result := make([]questView, len(quests))
		for i, q := range quests {
			result[i] = questToView(q)
		}
		c.JSON(http.StatusOK, gin.H{"quests": result})
	}
}

// getQuest returns a single quest by ID.
func getQuest(client *db.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quest id"})
			return
		}
		q, err := client.Quest.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "quest not found"})
			return
		}
		c.JSON(http.StatusOK, questToView(q))
	}
}