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

	r.GET("/hourly", getHourly)
	r.GET("/daily", getDaily)
	r.GET("/monthly", getMonthly)
	r.GET("/yearly", getYearly)

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

func getHourly(ctx *gin.Context) {
	date := ctx.Query("date")
	if date == "" {
		date = db.PDCollection.Meta.MinDate
	}

	m := db.PDCollection.FindOneByDate(date)
	ctx.JSON(http.StatusOK, m)
}

func getDaily(ctx *gin.Context) {
	date := ctx.Query("date")
	if date == "" {
		date = db.PDCollection.Meta.MinDate
	}

	m := db.PDCollection.FindDemandedDaily(date)
	ctx.JSON(http.StatusOK, m)
}

func getMonthly(ctx *gin.Context) {
	month := ctx.Query("month")
	if month == "" {
		month = db.PDCollection.Meta.MinDate[5:7]
	}
	year := ctx.Query("year")
	if year == "" {
		year = db.PDCollection.Meta.MinDate[0:4]
	}

	m := db.PDCollection.FindDemandedMonthly(year, month)
	ctx.JSON(http.StatusOK, m)
}

func getYearly(ctx *gin.Context) {
	year := ctx.Query("year")
	if year == "" {
		year = db.PDCollection.Meta.MinDate[0:4]
	}

	m := db.PDCollection.FindDemandedYearly(year)
	ctx.JSON(http.StatusOK, m)
}
