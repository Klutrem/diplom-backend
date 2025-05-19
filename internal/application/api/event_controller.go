package api

import (
	"main/internal/domain/events"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	logger       pkg.Logger
	eventService events.EventService
}

func NewEventController(logger pkg.Logger, eventService events.EventService) *EventController {
	return &EventController{
		logger:       logger,
		eventService: eventService,
	}
}

func (c *EventController) ListEvents(ctx *gin.Context) {
	namespace := ctx.Query("namespace")
	limitStr := ctx.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.logger.Errorf("invalid limit parameter: %s", limitStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	foundEvents, err := c.eventService.GetEvents(namespace, limit)
	if err != nil {
		c.logger.Errorf("failed to list events: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"events": foundEvents,
		"total":  len(foundEvents),
	})
}
