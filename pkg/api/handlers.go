package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Ivans's Internal Platform API " + Version,
	})
}

// Ready godoc
// @Summary      Check if the service is ready
// @Accept		 json
// @Produce      json
// @Success      200  {object}  interface{}
// @Failure      500  {object}	interface{}
// @Router       /healtz/ready [get]
func Ready(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}
