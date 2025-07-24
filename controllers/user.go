package controllers

import (
	"errors"
	"order-api/database"
	"order-api/models"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func CreateResponseUser(userModel models.User) User {
	return User{
		ID:        userModel.ID,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
	}
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.DB.Create(&user)
	response := CreateResponseUser(user)

	return c.Status(200).JSON(response)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}

	database.DB.Find(&users)
	responseUsers := []User{}

	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func findUser(id int, user *models.User) error {
	database.DB.Find(&user, id)
	if user.ID == 0 {
		return errors.New("User not found")
	}
	return nil
}
func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Id is not given or not numeric")
	}
	var user models.User

	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON("User not exist")
	}

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)

}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Id is not given or not numeric")
	}
	var user models.User

	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON("User not exist")
	}

	type updateUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	var updateData updateUser

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName

	database.DB.Save(&user)

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func RemoveUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Id is not given or not numeric")
	}
	var user models.User

	if err := findUser(id, &user); err != nil {
		return c.Status(404).JSON("User not exist")
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(500).SendString("failed to delete user")
	}

	return c.Status(200).JSON(fiber.Map{
		"msg": "Successfully deleted the user",
	})
}
