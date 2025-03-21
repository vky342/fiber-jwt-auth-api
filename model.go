package main

type User struct{
	Username string `json:"name"`
	Password string `json:"-"`    // prevetns marshlling to json
}

type FreshDonation struct {
	ID string `json:"id"`
	CreatedBy string `json:"name"`
	FoodName string `json:"foodname"`
	Quantity string `json:"quantity"`
	DateCreated string `json:"datecreated"`
	ExpiryDate string `json:"expirydate"`
	FoodPic []byte `json:"foodpic"`
}
