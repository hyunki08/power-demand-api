package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyunki08/power-demand-api/api/db"
)

func addPowerDemandRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/pd")

	// GET meta data
	r.GET("/", getMetadata)
	// GET one by date
	r.GET("/date", getByDate)
	// GET datas by range
	r.GET("/range", getByRange)
}

func getMetadata(ctx *gin.Context) {
	results := map[string]string{"minDate": db.PDCollection.Meta.MinDate, "maxDate": db.PDCollection.Meta.MaxDate}
	ctx.JSON(http.StatusOK, results)
}

func getByRange(ctx *gin.Context) {
	from := ctx.Query("from")
	if from == "" {
		from = db.PDCollection.Meta.MinDate
	}
	to := ctx.Query("to")
	if to == "" {
		to = db.PDCollection.Meta.MaxDate
	}

	ms := db.PDCollection.Find(from, to)
	ctx.JSON(http.StatusOK, ms)
}

func getByDate(ctx *gin.Context) {
	date := ctx.Query("date")
	if date == "" {
		date = db.PDCollection.Meta.MinDate
	}

	m := db.PDCollection.FindOneByDate(date)
	ctx.JSON(http.StatusOK, m)
}
