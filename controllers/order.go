package controllers

import (
	"errors"
	"order-api/database"
	"order-api/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

//ex:
// {
// 	id: 1,
// 	user: {
// 		id: 23,
// 		first_name: "some",
// 		last_name: "thing"
// 	},
// 	product:{
// 		id: 19,
// 		name: "Macbook",
// 		serial_number: "2345876"
// 	}
// }

type Order struct {
	ID        uint      `json:"id"`
	User      User      `json:"user"`
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"order_date"`
}

func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{
		ID:        order.ID,
		User:      user,
		Product:   product,
		CreatedAt: order.CreatedAt,
	}
}

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	respnseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(fiber.StatusOK).JSON(respnseOrder)
}

func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order
	if err := database.DB.Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	var responseOrders []Order
	for _, ord := range orders {
		var user models.User
		var product models.Product
		if err := database.DB.Find(&user, ord.UserRefer).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		if err := database.DB.Find(&product, ord.ProductRefer).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "product not found",
			})
		}
		responseUser := CreateResponseUser(user)
		responseProduct := CreateResponseProduct(product)
		responseOrders = append(responseOrders, CreateResponseOrder(ord, responseUser, responseProduct))
	}

	return c.Status(fiber.StatusOK).JSON(responseOrders)
}

func FindOrder(id int, order *models.Order) error {
	if err := database.DB.Find(&order, id).Error; err != nil {
		return errors.New("Order details not found")
	}
	return nil
}

func GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	var order models.Order
	if err := FindOrder(id, &order); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	var user models.User
	var product models.Product
	if err := database.DB.First(&user, order.UserRefer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}
	if err := database.DB.First(&product, order.ProductRefer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(err.Error())
	}

	responseUser := CreateResponseUser(user)
	responseProduct := CreateResponseProduct(product)
	responseOrder := CreateResponseOrder(order, responseUser, responseProduct)

	return c.Status(200).JSON(responseOrder)

}
