package api

import (
	"main/internal/infrastructure/kubernetes"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NamespaceController struct {
	kubernetesClient *kubernetes.KubernetesClient
}

func NewNamespaceController(kubernetesClient *kubernetes.KubernetesClient) *NamespaceController {
	return &NamespaceController{kubernetesClient: kubernetesClient}
}

func (c *NamespaceController) GetNamespaces(ctx *gin.Context) {
	namespaces, err := c.kubernetesClient.GetNamespaces()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, namespaces)
}
