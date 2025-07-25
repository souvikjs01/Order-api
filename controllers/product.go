package controllers

import (
	"errors"
	"order-api/database"
	"order-api/models"

	"github.com/gofiber/fiber/v2"
)

type Product struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) Product {
	return Product{
		ID:           productModel.ID,
		Name:         productModel.Name,
		SerialNumber: productModel.SerialNumber,
	}
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.DB.Create(&product)

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}

func GetProducts(c *fiber.Ctx) error {
	products := []models.Product{}

	database.DB.Find(&products)

	responseProducts := []Product{}

	for _, prod := range products {
		responseProducts = append(responseProducts, CreateResponseProduct(prod))
	}

	return c.Status(200).JSON(responseProducts)
}

func findProduct(id int, product *models.Product) error {
	database.DB.Find(&product, id)
	if product.ID == 0 {
		return errors.New("Product not exist")
	}
	return nil
}

func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Product id is not given")
	}

	var product models.Product

	if err := findProduct(id, &product); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Product id is not given")
	}

	var product models.Product
	if err := findProduct(id, &product); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	type updateProduct struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}
	var updateProductData updateProduct
	if err := c.BodyParser(&updateProductData); err != nil {
		c.Status(400).JSON("Invalid product")
	}

	product.Name = updateProductData.Name
	product.SerialNumber = updateProductData.SerialNumber
	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update data",
		})
	}

	responseProduct := CreateResponseProduct(product)

	return c.Status(fiber.StatusOK).JSON(responseProduct)
}

func DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Product id is not given")
	}

	var product models.Product
	if err := findProduct(id, &product); err != nil {
		return c.Status(404).JSON(err.Error())
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete the product",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"msg": "Product deleted successfully",
	})

}
