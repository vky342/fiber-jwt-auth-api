package main

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/dynamodb/v2"
	"golang.org/x/crypto/bcrypt"
)

func AuthHandlers(route fiber.Router, store dynamodb.Storage){

	route.Post("/register", func (c *fiber.Ctx) error {

		user := User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error" : "username and password required",
			})
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})	
		}

		user.Password = string(hashed)

		err = store.Set(user.Username, hashed, time.Second * 2)

		if err != nil{
			return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
				"error" : err.Error(),
			})
		}

		token,err := GenerateToken(&user)

		if err != nil{
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error" : err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name: "jwt",
			Value: token,
			HTTPOnly: !c.IsFromLocal(),
			Secure: !c.IsFromLocal(),
			MaxAge: 3600 * 24 * 7, // 7 days
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token" : token,
		})

	})


	route.Post("login", func(c *fiber.Ctx) error {
		user := User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error" : "username and password required",
			})
		}

		retrieved_password,err := store.Get(user.Username)

		if err != nil {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error" : "error fetching the password from DB",
				"mssg": "username wrong",
			})
		}

		//return c.SendString(string(retrieved_password) + " sent " + user.Password)

		err = bcrypt.CompareHashAndPassword(retrieved_password,[]byte(user.Password))

		if err != nil {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error" : "Password is wrong",
			})
		}

		token,err := GenerateToken(&user)

		if err != nil {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error" : "Token generation failed",
			})
		}
		
		c.Cookie(&fiber.Cookie{
			Name: "jwt",
			Value: token,
			HTTPOnly: !c.IsFromLocal(),
			Secure: !c.IsFromLocal(),
			MaxAge: 3600 * 24 * 7, // 7 days
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"msg" :"Login Succes",
			"token": token,
		})
	})

}