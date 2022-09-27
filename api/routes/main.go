package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run() {
	router.Use(CORSMiddleware())
	getRoutes()
	router.Run()
}

func getRoutes() {
	v1 := router.Group("/v1")
	addPowerDemandRoutes(v1)

	router.GET("/", getStatus)
}

func getStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ok")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
