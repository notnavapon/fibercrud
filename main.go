package main

import (
	"clean/infrastructure/http/router"
	"clean/infrastructure/initialize"
	"log"
)

func main() {
	config := initialize.Load()

	db := initialize.Database(config)
	deps := initialize.NewDependencies(db, config.JwtSecret)
	app := router.Router(deps, config.JwtSecret)

	log.Fatal(app.Listen(":" + config.ServerPort))
}
