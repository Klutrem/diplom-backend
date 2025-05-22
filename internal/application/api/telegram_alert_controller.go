package api

import (
	"main/internal/domain/alerts"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TelegramAlertController struct {
	logger  pkg.Logger
	service alerts.TelegramAlertService
}

func NewTelegramAlertController(logger pkg.Logger, service alerts.TelegramAlertService) *TelegramAlertController {
	return &TelegramAlertController{
		logger:  logger,
		service: service,
	}
}

func (c *TelegramAlertController) CreateAlert(ctx *gin.Context) {
	var alert alerts.TelegramAlert
	if err := ctx.ShouldBindJSON(&alert); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := c.service.CreateAlert(alert); err != nil {
		c.logger.Errorf("failed to create alert: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create alert"})
		return
	}

	response := alerts.TelegramAlertResponse{
		ID:        alert.ID,
		ChatID:    alert.ChatID,
		ThreadID:  alert.ThreadID,
		AlertType: alert.AlertType,
		Namespace: alert.Namespace,
		CreatedAt: alert.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, response)
}

func (c *TelegramAlertController) UpdateAlert(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	var alert alerts.TelegramAlert
	if err := ctx.ShouldBindJSON(&alert); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	alert.ID = id

	if err := c.service.UpdateAlert(alert); err != nil {
		c.logger.Errorf("failed to update alert: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update alert"})
		return
	}

	response := alerts.TelegramAlertResponse{
		ID:        alert.ID,
		ChatID:    alert.ChatID,
		ThreadID:  alert.ThreadID,
		AlertType: alert.AlertType,
		Namespace: alert.Namespace,
		CreatedAt: alert.CreatedAt,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *TelegramAlertController) DeleteAlert(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	if err := c.service.DeleteAlert(id); err != nil {
		c.logger.Errorf("failed to delete alert: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete alert"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "alert deleted successfully"})
}

func (c *TelegramAlertController) GetAlert(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	alert, err := c.service.GetAlert(id)
	if err != nil {
		c.logger.Errorf("failed to get alert: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get alert"})
		return
	}

	response := alerts.TelegramAlertResponse{
		ID:        alert.ID,
		ChatID:    alert.ChatID,
		ThreadID:  alert.ThreadID,
		AlertType: alert.AlertType,
		Namespace: alert.Namespace,
		CreatedAt: alert.CreatedAt,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *TelegramAlertController) GetAlertsByNamespace(ctx *gin.Context) {
	namespace := ctx.Param("namespace")
	if namespace == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "namespace parameter is required"})
		return
	}

	resp, err := c.service.GetAlertsByNamespace(namespace)
	if err != nil {
		c.logger.Errorf("failed to get alerts by namespace: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get alerts"})
		return
	}

	responses := make([]alerts.TelegramAlertResponse, len(resp))
	for i, alert := range resp {
		responses[i] = alerts.TelegramAlertResponse{
			ID:        alert.ID,
			ChatID:    alert.ChatID,
			ThreadID:  alert.ThreadID,
			AlertType: alert.AlertType,
			Namespace: alert.Namespace,
			CreatedAt: alert.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"alerts": responses,
		"total":  len(responses),
	})
}

func (c *TelegramAlertController) GetAllAlerts(ctx *gin.Context) {
	resp, err := c.service.GetAllAlerts()
	if err != nil {
		c.logger.Errorf("failed to get all alerts: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get alerts"})
		return
	}

	responses := make([]alerts.TelegramAlertResponse, len(resp))
	for i, alert := range resp {
		responses[i] = alerts.TelegramAlertResponse{
			ID:        alert.ID,
			ChatID:    alert.ChatID,
			ThreadID:  alert.ThreadID,
			AlertType: alert.AlertType,
			Namespace: alert.Namespace,
			CreatedAt: alert.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"alerts": responses,
		"total":  len(responses),
	})
}
