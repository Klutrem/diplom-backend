package api

import (
	"main/internal/infrastructure/kubernetes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PodController struct {
	kubernetesClient *kubernetes.KubernetesClient
}

func NewPodController(kubernetesClient *kubernetes.KubernetesClient) *PodController {
	return &PodController{kubernetesClient: kubernetesClient}
}

func (c *PodController) GetPods(ctx *gin.Context) {
	namespace := ctx.DefaultQuery("namespace", "default")
	println(namespace)
	pods, err := c.kubernetesClient.GetPods(namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pods)
}

func (c *PodController) GetPodMetrics(ctx *gin.Context) {
	namespace := ctx.Param("namespace")
	podName := ctx.Param("pod")

	startStr := ctx.DefaultQuery("start", time.Now().Add(-1*time.Hour).Format(time.RFC3339))
	endStr := ctx.DefaultQuery("end", time.Now().Format(time.RFC3339))
	stepStr := ctx.DefaultQuery("step", "15s")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time format"})
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time format"})
		return
	}

	step, err := time.ParseDuration(stepStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid step duration"})
		return
	}

	metrics, err := c.kubernetesClient.GetPodHistoricalMetrics(namespace, podName, start, end, step)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, metrics)
}
