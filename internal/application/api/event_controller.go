package api

import (
	"main/internal/domain/events"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EventController handles HTTP requests for events
type EventController struct {
	logger       pkg.Logger
	eventService *events.EventService
}

// NewEventController creates a new EventController
func NewEventController(logger pkg.Logger, eventService *events.EventService) *EventController {
	return &EventController{
		logger:       logger,
		eventService: eventService,
	}
}

// ListEvents handles GET /events
func (c *EventController) ListEvents(ctx *gin.Context) {
	namespace := ctx.Query("namespace")
	limitStr := ctx.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.logger.Errorf("invalid limit parameter: %s", limitStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	events, err := c.eventService.GetEvents(ctx, namespace, limit)
	if err != nil {
		c.logger.Errorf("failed to list events: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"events": events,
		"total":  len(events),
	})
}
