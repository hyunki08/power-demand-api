package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addPowerDemandRoutes(rg *gin.RouterGroup) {
	ping := rg.Group("/pd")

	ping.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello")
	})
}
