package prxmx

import (
	"i2/pkg/types"

	"github.com/gin-gonic/gin"
)

func AddRoutes(api *gin.RouterGroup, config *types.Config) {
	cluster := NewCluster(config.Proxmox.URL, config.Proxmox.User, config.Proxmox.Pass)
	api.GET("/proxmox/nodes", cluster.handlerGetClusterNodes)
	api.GET("/proxmox/vms", cluster.handlerGetVirtualMachines)
}
