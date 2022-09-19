package routes

import "github.com/gin-gonic/gin"

var router = gin.Default()

func Run() {
	getRoutes()
	router.Run()
}

func getRoutes() {
	v1 := router.Group("/v1")
	addMetaRoutes(v1)
	addPowerDemandRoutes(v1)
}
