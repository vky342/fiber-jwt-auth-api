package main

import (
	"encoding/json"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/dynamodb/v2"
	"github.com/google/uuid"
)

func FreshDonationHandlers(route fiber.Router, store dynamodb.Storage){


	// create new donation
	route.Post("/create", func (c *fiber.Ctx) error  {

		user := new(User)
	
		token, err :=AuthMiddleware(c)
		if err != nil{
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"error": "Auth Middle-ware failure RFD-30",
			})
		}
	
	
		user,err = ExtractUserFromToken(token)
		if err != nil {
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error": "User Credentials extraction failed",
			})
		}
			
	
		newDonation := FreshDonation{
			ID: uuid.New().String(),
			CreatedBy: user.Username,
			FoodName: c.FormValue("foodname"),
			Quantity: c.FormValue("quantity"),
			DateCreated: c.FormValue("datecreated"),
			ExpiryDate: c.FormValue("expirydate"),
			FoodPic: []byte(c.FormValue("foodpic")),
		}
	
		byte_array,err := json.Marshal(newDonation)
	
		if err != nil{
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error": "Json Marshaling failed RFD-59",
			})
		}
	
		err = store.Set(newDonation.ID,byte_array,time.Second * 3)
	
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Writing into DB failed",
			})
		}
	
		return c.SendString("created new Donation!")
	
	})



	
	
	// Get all donations
	route.Get("/donations", func(c *fiber.Ctx) error {

		donations := make([]FreshDonation, 0)

		Keys,err := listPrimaryKeys(&store,"Fresh_Donations","k")

		if err != nil{
			return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{
				"error": "Listing Primary Key failure RFD-23",
			})
		}

		for _, key := range Keys { 
		
			data, err := store.Get(key)
			if err != nil {
				continue // Skip failed reads
			}

			var donation FreshDonation
			err = json.Unmarshal(data, &donation)
			if err != nil {
				continue // skip corrupted data
			}

			donations = append(donations, donation)
		}

		return c.JSON(donations)
	})



	// Get available donations only


	// Get available Donations by user
}