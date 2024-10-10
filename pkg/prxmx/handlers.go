package prxmx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetClusterNodes godoc
// @Summary Get cluster nodes
// @Description Get cluster nodes
// @Tags proxmox
// @Accept json
// @Produce json
// @Success 200 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /proxmox/nodes [get]
func (cluster *Cluster) handlerGetClusterNodes(c *gin.Context) {
	nodes, err := cluster.Client.Nodes(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

// GetVirtualMachines godoc
// @Summary Get virtual machines
// @Description Get virtual machines
// @Tags proxmox
// @Accept json
// @Produce json
// @Success 200 {array} []Node
// @Failure 500 {object} interface{}
// @Router /proxmox/vms [get]
func (cluster *Cluster) handlerGetVirtualMachines(c *gin.Context) {
	nodes, err := cluster.GetVMs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}
