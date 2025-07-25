package routes

import (
	"order-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(r *fiber.App) {
	// user routes
	r.Post("/user", controllers.CreateUser)
	r.Get("/users", controllers.GetUsers)
	r.Get("/user/:id", controllers.GetUser)
	r.Put("/user/update/:id", controllers.UpdateUser)
	r.Delete("/user/remove/:id", controllers.RemoveUser)

	// order routes
	r.Post("/product", controllers.CreateProduct)
	r.Get("/product/all", controllers.GetProducts)
	r.Get("/product/:id", controllers.GetProduct)
	r.Put("/product/update/:id", controllers.UpdateProduct)
	r.Delete("/product/remove/:id", controllers.DeleteProduct)
}
