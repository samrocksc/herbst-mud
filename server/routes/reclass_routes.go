package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"herbst-server/service"
)

// RegisterReclassRoutes registers reclass/rerace and history endpoints.
func RegisterReclassRoutes(r *gin.Engine, svc *service.Container) {
	characters := r.Group("/api/characters")
	{
		characters.POST("/:id/reclass", reclassHandler(svc))
		characters.POST("/:id/rerace", reraceHandler(svc))
		characters.GET("/:id/class-history", classHistoryHandler(svc))
		characters.GET("/:id/race-history", raceHistoryHandler(svc))
	}
}

// reclassRequest is the body for POST /api/characters/:id/reclass
type reclassRequest struct {
	FactionID int `json:"faction_id"`
}

func reclassHandler(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}

		var req reclassRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if req.FactionID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "faction_id is required"})
			return
		}

		err = svc.ReclassRerace.Reclass(c.Request.Context(), charID, req.FactionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "reclass successful"})
	}
}

// reraceRequest is the body for POST /api/characters/:id/rerace
type reraceRequest struct {
	RaceName string `json:"race_name"`
}

func reraceHandler(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}

		var req reraceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if req.RaceName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "race_name is required"})
			return
		}

		err = svc.ReclassRerace.Rerace(c.Request.Context(), charID, req.RaceName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "rerace successful"})
	}
}

func classHistoryHandler(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}

		history, err := svc.ReclassRerace.GetClassHistory(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"history": history})
	}
}

func raceHistoryHandler(svc *service.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		charID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid character id"})
			return
		}

		history, err := svc.ReclassRerace.GetRaceHistory(c.Request.Context(), charID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"history": history})
	}
}