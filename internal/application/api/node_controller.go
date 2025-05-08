package api

import (
	"main/internal/infrastructure/kubernetes"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NodeController struct {
	kubernetesClient *kubernetes.KubernetesClient
}

func NewNodeController(kubernetesClient *kubernetes.KubernetesClient) *NodeController {
	return &NodeController{kubernetesClient: kubernetesClient}
}

func (c *NodeController) GetNodes(ctx *gin.Context) {
	nodes, err := c.kubernetesClient.GetNodes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, nodes)
}
