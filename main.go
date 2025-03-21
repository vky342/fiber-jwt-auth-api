package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/dynamodb/v2"
)

func main() {

	//create instances here
	app := fiber.New()

	store_restaurants := InitializeDB(dynamodb.Config{Endpoint: "http://localhost:8000",Table: "Restaurants"})

	store_ngos := InitializeDB(dynamodb.Config{Endpoint: "http://localhost:8000",Table: "Ngos"})

	store_fresh_Donation := InitializeDB(dynamodb.Config{Endpoint: "http://localhost:8000", Table: "Fresh_Donations"})


	//Handlers here
	AuthHandlers(app.Group("/restaurants/auth"),store_restaurants)

	AuthHandlers(app.Group("/ngos/auth"),store_ngos)

	FreshDonationHandlers(app.Group("/freshdonations"),store_fresh_Donation)


	//start the server on port 3000
	app.Listen(":3000")

}