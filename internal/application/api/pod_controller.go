package api

import (
	"main/internal/infrastructure/kubernetes"
	"net/http"

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
