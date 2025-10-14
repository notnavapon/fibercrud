package main

import (
	"clean/internal/initialize"
	"log"
)

func main() {
	config := initialize.Load()

	db := initialize.Database(config)
	deps := initialize.NewDependencies(db, config.JwtSecret)
	app := initialize.Router(deps, config.JwtSecret)

	log.Fatal(app.Listen(":" + config.ServerPort))
}
