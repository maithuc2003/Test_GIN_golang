package main

import (
	"github.com/maithuc2003/GIN_golang_framework/config"
	"github.com/maithuc2003/GIN_golang_framework/internal/routes"
)

func main() {
	db := config.ConnectDB()
	r := routes.SetupRouter(db)
	r.Run(":8080")
}
