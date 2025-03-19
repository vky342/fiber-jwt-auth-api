package main

type User struct{
	Username string `json:"name"`
	Password string `json:"-"`    // prevetns marshlling to json
}