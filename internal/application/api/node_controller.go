package api

import (
	"main/internal/infrastructure/kubernetes"
	"main/pkg"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NodeController struct {
	kubernetesClient *kubernetes.KubernetesClient
	logger           pkg.Logger
}

func NewNodeController(kubernetesClient *kubernetes.KubernetesClient) *NodeController {
	return &NodeController{kubernetesClient: kubernetesClient}
}

func (c *NodeController) GetNodes(ctx *gin.Context) {
	nodes, err := c.kubernetesClient.GetNodes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.logger.Errorf("Failed to get nodes: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, nodes)
}

func (c *NodeController) GetNodeMetrics(ctx *gin.Context) {
	nodeName := ctx.Param("node")

	startStr := ctx.DefaultQuery("start", time.Now().Add(-1*time.Hour).Format(time.RFC3339))
	endStr := ctx.DefaultQuery("end", time.Now().Format(time.RFC3339))
	stepStr := ctx.DefaultQuery("step", "15s")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time format"})
		c.logger.Errorf("Failed to parse start time: %v", err)
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time format"})
		c.logger.Errorf("Failed to parse end time: %v", err)
		return
	}

	step, err := time.ParseDuration(stepStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid step duration"})
		c.logger.Errorf("Failed to parse step duration: %v", err)
		return
	}

	metrics, err := c.kubernetesClient.GetNodeHistoricalMetrics(nodeName, start, end, step)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.logger.Errorf("Failed to get node metrics: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, metrics)
}
