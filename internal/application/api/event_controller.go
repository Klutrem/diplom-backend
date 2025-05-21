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
	namespace := ctx.DefaultQuery("namespace", "default")
	eventType := ctx.Query("type")
	limitStr := ctx.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.logger.Errorf("invalid limit parameter: %s", limitStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	foundEvents, err := c.eventService.GetEvents(namespace, eventType, limit)
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

func (c *EventController) GetWatchedNamespaces(ctx *gin.Context) {
	namespaces, err := c.eventService.GetWatchedNamespaces()
	if err != nil {
		c.logger.Errorf("failed to get watched namespaces: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get watched namespaces"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"namespaces": namespaces,
		"total":      len(namespaces),
	})
}

func (c *EventController) AddWatchedNamespace(ctx *gin.Context) {
	namespace := ctx.Query("namespace")

	if err := c.eventService.AddNamespaceToWatch(ctx, namespace); err != nil {
		c.logger.Errorf("failed to add namespace to watch: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add namespace to watch"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "namespace added to watch list"})
}

func (c *EventController) RemoveWatchedNamespace(ctx *gin.Context) {
	namespace := ctx.Param("namespace")
	if namespace == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "namespace parameter is required"})
		return
	}

	if err := c.eventService.RemoveNamespaceFromWatch(namespace); err != nil {
		c.logger.Errorf("failed to remove namespace from watch: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove namespace from watch"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "namespace removed from watch list"})
}
