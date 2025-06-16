package main

import (
	"github.com/maithuc2003/Test_GIN_golang/config"
	"github.com/maithuc2003/Test_GIN_golang/internal/routes"
)

func main() {
	db := config.ConnectDB()
	r := routes.SetupRouter(db)
	r.Run(":8080")
}
