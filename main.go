package main

import (
	"log"
	"order-api/database"
	"order-api/routes"

	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to my order api")
}
func main() {
	database.ConnectDB()
	app := fiber.New()
	app.Get("/", welcome)

	routes.Router(app)

	log.Fatal(app.Listen(":3000"))
}
