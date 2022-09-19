package main

import (
	"github.com/hyunki08/power-demand-api/api/db"
	"github.com/hyunki08/power-demand-api/api/routes"
)

func main() {
	db.Run()
	defer db.Disconnect()
	routes.Run()
}
