package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run() {
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
